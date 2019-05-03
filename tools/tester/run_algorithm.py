#!/usr/bin/env python3

'''
Script used for running all instances in a given algorithm.
'''

import argparse
import subprocess
import os.path
import shutil


def main():
    parser = argparse.ArgumentParser(
        description='Run all instances of a given algorithm.')
    parser.add_argument(
        '--algorithm', '-a', help="The algorithm to run all instances for.")
    args = parser.parse_args()

    # create results directory
    dirname = f'{args.algorithm.upper()}_RESULTS/'
    if not os.path.exists(dirname):
        print(f'Making directory: {dirname}')
        os.mkdir(dirname)
    else:  # remove and recreate if it already exists
        print(f'Removing directory: {dirname}')
        shutil.rmtree(dirname)
        print(f'Making directory: {dirname}')
        os.mkdir(dirname)

    # iterate through item sizes
    for item_size_percent in [25, 50, 75]:
        # iterate through item counts
        for item_count in [50, 100, 500]:
            # iterate through item variances
            for item_variance in [1, 2, 3]:
                # create 10,000 instances per test case
                for i in range(10000):
                    filename = f'binpacking{i}_{item_count}count_{100}max_{item_size_percent}center_{item_variance}variability_{args.algorithm}'
                    input_json = f'{args.algorithm}/{filename}.json'
                    output_json = f'{dirname}/{filename}_results.json'

                    bashCommand = f'./tester -file={input_json} -output={output_json}'
                    print(f'Running bash command: {bashCommand}')
                    process = subprocess.Popen(
                        bashCommand.split(), stdout=subprocess.PIPE)
                    output, error = process.communicate()


if __name__ == '__main__':
    main()
