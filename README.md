# Assignment 2: Genome Assembly and Evaluation

TODO:
- [ ] Add brief introduction
- [ ] Add executable scripts 
- [ ] Add example output
- [ ] Add instructions on how to run the scripts

for the pdf report:
- [ ] brief description of each algorithm
- [ ] assembly results and analysis
- [ ] visualization figures and comparisons
- [ ] evaluation metrics with interpretation
- [ ] discussion of challenges and solutions applied including runtime and memory usage

## LLM usage

- I used Cursor IDE to write the code for this assignment.
  - Cursor has Cursor Tab functionality that auto completes lines and functions.
  - I used it while coding for all files.
  - It uses a mix of LLMs including GPT-4, Claude, Deepseek and Gemini.
  - I used the pro tier of Cursor.
- I used ChatGPT with the `gpt-4o` model to write:
  - test cases in `dbg_test.go` file.
  - `exportToGFA` function in `dbg.go` file.
  - Also, for help with running bio-tools specified in the assignment.
- I tried using Aider with Claude 3.7 and gemini 2.5 flash to write functions in `dbg.go` file but all attempts failed.


## Running the code

### Task 1.3.1:
- `go run src/dbg.go toy_dataset/reads_b.fastq 40`
  - this will generate `reads_b_dbg_k_40.fasta` and `reads_b_dbg_k_40.gfa` files in the `toy_dataset` directory.
  - use Bandage to visualize the graph in `reads_b_dbg_k_40.gfa` file.

### Task 1.3.2:
1. `go run src/dbg.go toy_dataset/reads_r.fastq 35`
2. `go run src/dbg.go toy_dataset/reads_r.fastq 45`
  - this will generate `reads_r_dbg_k_35.fasta`, `reads_r_dbg_k_35.gfa`, `reads_r_dbg_k_45.fasta`, `reads_r_dbg_k_45.gfa` files in the `toy_dataset` directory.
  - use Bandage to visualize the graphs in `reads_r_dbg_k_35.gfa` and `reads_r_dbg_k_45.gfa` files.
3. Then, inside the `toy_dataset` directory, run `quast.py -r reference_r.fasta -o quast_output_r_k_35 reads_r_dbg_k_35.fasta` and `quast.py -r reference_r.fasta -o quast_output_r_k_45 reads_r_dbg_k_45.fasta` to get the assembly metrics.

### Task 1.3.3:
Similar to Task 1.3.2, run the following commands to get `.fasta` and `.gfa` files for `synthetic_dataset/reads/` directory.
```bash
go run src/dbg.go synthetic_dataset/reads/no_error_reads_hiseq_5k.fastq 40
go run src/dbg.go synthetic_dataset/reads/no_error_ont_hq_50x.fastq 40
```

TODOs:
- fix code for error reads.
- add `olc.go` script runs.
- add quast commands to get assembly metrics.
