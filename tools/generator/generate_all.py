#!/usr/bin/env python3
import subprocess, shutil, os, os.path

''' 
Python script used for generating all of the 
JSON files needed for the test plan.
'''

# algorithms in use
ALGORITHMS = [
    "NextFit",
    "FirstFit",
    "FirstFitDecreasing",
    "BestFit",
    "BestFitDecreasing",
    "PackingConstraint",
    "BinCompletion",
    "ModifiedFirstFitDecreasing"
]

DUPLICATES = 10000

def main():
    # iterate through algorithms
    for algorithm in ALGORITHMS:
        if os.path.isfile(f'{algorithm}.tar'):
            os.remove(f'{algorithm}.tar')


        # iterate through item counts
        for item_size_percent in [25, 50, 75]:
            # iterate through item variances
            for item_count in [50, 100, 500]:
                # create 10,000 instances per test case 
                for item_variance in [1, 2, 3]:
                    # generate test case
                        bashCommand = f'./generator -algorithm={algorithm} -count={item_count} -dups={DUPLICATES} -variability={item_variance} -center={item_size_percent} -output={algorithm}'
                        print(f'Running bash command: {bashCommand}')
                        process = subprocess.Popen(bashCommand.split(), stdout=subprocess.PIPE)
                        output, error = process.communicate()

        # print(f'Making tar file: {algorithm}.tar')
        # shutil.make_archive(algorithm, 'tar', algorithm)
        # print(f'Removing directory: {algorithm}')
        # shutil.rmtree(algorithm)

if __name__ == '__main__':
    main()
