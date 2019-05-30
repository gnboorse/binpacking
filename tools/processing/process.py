#!/usr/bin/env python3
import sqlalchemy
import tarfile
import argparse
import os
import os.path
import json
from tqdm import tqdm
from sqlalchemy import Table, Column, Integer, String, MetaData, Float
import shutil


def get_digits(s):
    return ''.join(c for c in s if c.isdigit())


def parse_test_case_name(tc_name):
    """
    Method used to get data from the test case name.
    """
    split_tc = tc_name.split(
        '_')  # only interested in indices 1, 3, and 4 in this split
    return {
        'count': get_digits(split_tc[1]),
        'center': get_digits(split_tc[3]),
        'variability': get_digits(split_tc[4])
    }


def main():
    parser = argparse.ArgumentParser(
        description='utility for processing bin packing data')
    parser.add_argument('--resultsTar', '-r', required=True,
                        help='tar file containing results JSON files')
    parser.add_argument('--testCasesTar', '-t', required=True,
                        help='tar file continaing test case descriptions')
    parser.add_argument('--output', '-o', required=True,
                        help='filename of output db')
    parser.add_argument(
        '--algorithm', '-a', required=True, help="algorithm to process files for")
    args = parser.parse_args()

    print('Extracting results tar: %s' % args.resultsTar)
    tar = tarfile.open(args.resultsTar)
    tar.extractall()
    tar.close()

    print('Extracting test cases tar: %s' % args.testCasesTar)
    tar = tarfile.open(args.testCasesTar)
    tar.extractall()
    tar.close()

    test_cases_dir = args.algorithm
    results_dir = args.algorithm.upper() + "_RESULTS"
    if not os.path.exists(test_cases_dir):
        print(f'Test cases path does not exist: {test_cases_dir}')
        return
    
    if not os.path.exists(results_dir):
        print(f'Results path does not exist: {results_dir}')
        return

    print(f'Creating database connection to: {args.output}')
    sql_engine = sqlalchemy.create_engine(f'sqlite:///{args.output}')
    metadata = MetaData(sql_engine)
    bin_packing_results_table = Table('bin_packing_results', metadata,
                                      Column('Id', Integer,
                                             primary_key=True, nullable=False),
                                      Column('count', Integer),
                                      Column('center', Integer),
                                      Column('variability', Integer),
                                      Column('lower_bound', Integer),
                                      Column('algorithm', Integer),
                                      Column('solution_bin_count', Integer),
                                      Column('solution_time', Integer),
                                      Column('solution_optimality', Float))

    metadata.create_all(sql_engine)
    conn = sql_engine.connect()

    batch_size = 20 # size of the batches to insert 
    batch_holder = []
    # iterate over all test cases
    for test_case_file in tqdm(os.listdir(test_cases_dir)):
        try:
            test_case_name = os.path.splitext(test_case_file)[0]
            results_file = f'{test_case_name}_results.json'

            test_case_data = None
            test_results_data = None
            if os.path.exists(os.path.join(results_dir, results_file)):
                with open(os.path.join(test_cases_dir, test_case_file)) as fh:
                    test_case_data = json.load(fh)
                with open(os.path.join(results_dir, results_file)) as fh:
                    test_results_data = json.load(fh)

                data_dict = parse_test_case_name(test_case_name)
                data_dict['lower_bound'] = test_case_data['lowerBound']
                data_dict['algorithm'] = test_case_data['algorithm']
                data_dict['solution_bin_count'] = test_results_data['count']
                data_dict['solution_time'] = test_results_data['solution_time']
                data_dict['solution_optimality'] = round(float(data_dict['solution_bin_count']) / float(data_dict['lower_bound']), ndigits=5)

                # add data to batch
                batch_holder.append(data_dict)

                # data dict now contains all data we care about in this test case
                if len(batch_holder) >= batch_size:
                    # time to insert
                    conn.execute(bin_packing_results_table.insert(), batch_holder)
                    batch_holder = []
        except Exception:
            continue

    print(f'Removing test cases dir: {test_cases_dir}')
    shutil.rmtree(test_cases_dir)

    print(f'Removing results dir: {results_dir}')
    shutil.rmtree(results_dir)


if __name__ == '__main__':
    main()
