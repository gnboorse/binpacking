#!/usr/bin/env python3

'''
Script used for running all instances in a given directory of test cases.
'''
import os.path
import argparse
import tqdm
import shutil
import subprocess


def main():
    parser = argparse.ArgumentParser(
        description='Run all instances of a given algorithm.')
    parser.add_argument(
        '--file', '-f', help="The directory containing files.")
    args = parser.parse_args()

    # create results directory
    results_dirname = f'{args.file.upper()}_RESULTS/'
    if not os.path.exists(results_dirname):
        print(f'Making directory: {results_dirname}')
        os.mkdir(results_dirname)
    else:  # remove and recreate if it already exists
        print(f'Removing directory: {results_dirname}')
        shutil.rmtree(results_dirname)
        print(f'Making directory: {results_dirname}')
        os.mkdir(results_dirname)

    for input_json in tqdm(os.listdir(args.file)):
        output_json_filename = f'{os.path.splitext(input_json)[0]}_results.json'
        bashCommand = f'./tester -file={os.path.join(args.file,input_json)} -output={os.path.join(results_dirname, output_json_filename)}'
        # print(f'Running bash command: {bashCommand}')
        process = subprocess.Popen(
            bashCommand.split(), stdout=subprocess.PIPE)
        output, error = process.communicate()


if __name__ == '__main__':
    main()
