INSTRUCTIONS from CREATOR:
1. clone
2. in the directory togzhan-project execute command: go install counter
3. and when it's done, execute command: counter
4. answer question about number of go routines
5. then answer question about regeneration of file (EXPECT 30 seconds delay if regeneration was chosen, it's safe mechanism to wait until file is finished being created before starting to count numbers in it)
6. if you did not want to regenerate sample_file_10k.json will be used (if you want to check code with 1,000,000 objects, you have to regenerate file but then change 3rd parameter in line 71 to 1000000 in the main.go file)
7. done, you should see sum and time it took to count it
