# Assignment 2: Genome Assembly and Evaluation

- [Assignment 2: Genome Assembly and Evaluation](#assignment-2-genome-assembly-and-evaluation)
  - [Introduction](#introduction)
  - [Running the code](#running-the-code)
    - [Task 1.3.1 DBG Assembly Graph on reads\_b](#task-131-dbg-assembly-graph-on-reads_b)
    - [Task 1.3.2 DBG on reads\_r](#task-132-dbg-on-reads_r)
    - [Task 1.3.3 DBG and OLC on MERS virus data](#task-133-dbg-and-olc-on-mers-virus-data)
    - [Task 1.3.4 SPAdes](#task-134-spades)
    - [Task 2.1 Genome Assembly](#task-21-genome-assembly)
    - [Task 2.2 Evaluation Scripts](#task-22-evaluation-scripts)
      - [QUAST](#quast)
      - [BUSCO](#busco)
      - [Merqury](#merqury)
      - [Inspector](#inspector)
  - [LLM usage](#llm-usage)


## Introduction

In this assignment, I had the opportunity to work with assembling and evaluating genome assemblies. 

In the first part of the assignment, I implemented the De Bruijn Graph (DBG) and Overlap-Layout-Consensus (OLC) algorithms.
The assignment sheet provided a toy dataset to test the algorithms and a larger more realistic dataset of Middle East respiratory
syndrome-related coronavirus (MERS-CoV) genome. I evaluated the performance of these algorithms on the datasets and compared them against each other and against a more established tool, SPAdes.

In the second part of the assignment, I worked with a Scincus mitranus lizard genome and evaluated the performance of the assembly using the tools described in the assignment sheet. 

Analysis of the results are provided in the report found in the `report/` directory. The code is available in the `src/` directory. `toy_dataset/` and `synthetic_dataset/` contain resulting files from the algorithms as described below.

## Running the code

`src/` directory contains the source code for DBG and OLC implementations written in Go. These files are only needed for Task 1 to run the DBG and OLC algorithms.
```bash
$ tree src
src
├── dbg.go
├── main.go
└── olc.go

1 directory, 3 files
```

To run the code, use the following commands:
```bash
go run . dbg <reads_file> <k> [gfa]
go run . olc <reads_file> <min_overlap>
```
`gfa` is optional and if provided, the program will generate a GFA file for the assembly graph. The output files will be saved in the same path as the input file with `_dbg_k_<k>.fasta` and `_dbg_k_<k>.gfa` extensions or `_olc_<min_overlap>.fasta` extension.

### Task 1.3.1 DBG Assembly Graph on reads_b
- `go run . dbg toy_dataset/reads_b.fastq 40 gfa`
  - this will generate `reads_b_dbg_k_40.fasta` and `reads_b_dbg_k_40.gfa` files in the `toy_dataset` directory.
  - use Bandage to visualize the graph in `reads_b_dbg_k_40.gfa` file.

### Task 1.3.2 DBG on reads_r
1. `go run . dbg toy_dataset/reads_r.fastq 35 gfa`
2. `go run . dbg toy_dataset/reads_r.fastq 45 gfa`
  - this will generate `reads_r_dbg_k_35.fasta`, `reads_r_dbg_k_35.gfa`, `reads_r_dbg_k_45.fasta`, `reads_r_dbg_k_45.gfa` files in the `toy_dataset` directory.
  - use Bandage to visualize the graphs in `reads_r_dbg_k_35.gfa` and `reads_r_dbg_k_45.gfa` files.
3. Then, inside the `toy_dataset` directory, run `quast.py -r reference_r.fasta -o quast_output_r_k_35 reads_r_dbg_k_35.fasta` and `quast.py -r reference_r.fasta -o quast_output_r_k_45 reads_r_dbg_k_45.fasta` to get the assembly metrics.

### Task 1.3.3 DBG and OLC on MERS virus data
Similar to Task 1.3.2, run the following commands to get `.fasta` files for `synthetic_dataset/reads/` directory.
```bash
go run . dbg synthetic_dataset/reads/no_error_reads_hiseq_5k.fastq 40
go run . dbg synthetic_dataset/reads/no_error_ont_hq_50x.fastq 40

go run . olc synthetic_dataset/reads/no_error_reads_hiseq_5k.fastq 40
go run . olc synthetic_dataset/reads/no_error_ont_hq_50x.fastq 40
```

Then, run the following commands to get the QUAST evaluation results for the no-error reads.
```bash
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_reads_hiseq_5k_dbg_k_40 reads/no_error_reads_hiseq_5k_dbg_k_40.fasta
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_ont_hq_50x_dbg_k_40 reads/no_error_ont_hq_50x_dbg_k_40.fasta

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_reads_hiseq_5k_olc_40 reads/no_error_reads_hiseq_5k_olc_40.fasta
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_ont_hq_50x_olc_40 reads/no_error_ont_hq_50x_olc_40.fasta
```

Run the following commands to get the QUAST evaluation results for the reads with errors.
```bash
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_ont_hq_50x_dbg_k_40 reads/ont_hq_50x_dbg_k_40.fasta
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_reads_hiseq_5k_dbg_k_40 reads/reads_hiseq_5k_dbg_k_40.fasta

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_ont_hq_50x_olc_40 reads/ont_hq_50x_olc_40.fasta
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_reads_hiseq_5k_olc_40 reads/reads_hiseq_5k_olc_40.fasta
```

Results should be available then in the `synthetic_dataset/` directory with the `quast_` prefix. A readable report should be available in pdf format in each of the generated directories.

### Task 1.3.4 SPAdes
- Download SPAdes 4.1.0 (I couldn't get it from Ibex modules):
```bash
wget https://github.com/ablab/spades/releases/download/v4.1.0/SPAdes-4.1.0-Linux.tar.gz
 tar -xzf SPAdes-4.1.0-Linux.tar.gz
```

Run the following commands to get the SPAdes assembly results for no-error and error reads.
```bash
SPAdes-4.1.0-Linux/bin/spades.py --phred-offset 33 -s reads/no_error_reads_hiseq_5k.fastq -o spades_output_no_error_reads_hiseq_5k
SPAdes-4.1.0-Linux/bin/spades.py --phred-offset 33 -s reads/reads_hiseq_5k.fastq -o spades_output_hiseq_5k

SPAdes-4.1.0-Linux/bin/spades.py --phred-offset 33 -s reads/no_error_ont_hq_50x.fastq -o spades_output_no_error_ont_hq_50x
SPAdes-4.1.0-Linux/bin/spades.py --phred-offset 33 -s reads/ont_hq_50x.fastq -o spades_output_ont_hq_50x
```

Run the following commands to get the QUAST evaluation results for the SPAdes assemblies.
```bash
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_reads_hiseq_5k_spades spades_output_no_error_reads_hiseq_5k/contigs.fasta 

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_reads_hiseq_5k_spades spades_output_hiseq_5k/contigs.fasta 

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_ont_hq_50x_spades spades_output_no_error_ont_hq_50x/contigs.fasta 

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_ont_hq_50x_spades spades_output_ont_hq_50x/contigs.fasta

```

Similarly, pdf reports should be available in the generated directories with the `quast_` prefix.

### Task 2.1 Genome Assembly

`asm_scinmit/` directory is used to store the data and results for the lizard genome assembly.

Assembly plan:
```bash
# 0- Allocate resources with slurm script (part of the head of the script):
#SBATCH --cpus-per-task=64
#SBATCH --mem=256G
#SBATCH --account=cs249

# 1- Setup
module load gcc/12.2.0 hifiasm/0.19 minimap2 samtools seqtk
mkdir -p asm_scinmit/{assembly,data,eval,tmp}
cd asm_scinmit

# 2- Fetch data
export RAW=/ibex/reference/course/cs249/lizard/input
# Create symbolic links inside project tree
ln -s $RAW/pacbio/*.fastq.gz data/pacbio/
ln -s $RAW/hic/*.fastq.gz data/hic/
ln -s $RAW/ont/*.fastq.gz data/ont/

# 3- Assemble primary contigs (took around 9 hours and 16 minutes)
# Real time: 33467.576 sec; CPU: 1523075.897 sec; Peak RSS: 170.915 GB
hifiasm -o assembly/scinmit_hicul \
        -t64 --primary -D10 \
        --h1 data/hic/lizard_hic_R1.fastq.gz \
        --h2 data/hic/lizard_hic_R2.fastq.gz \
        --ul data/ont/lizard_ont.fastq.gz \
        data/pacbio/lizard_liver_seq.fastq.gz \
        data/pacbio/lizard_liver_rev.fastq.gz

# convert GFA to FASTA
awk '$1=="S"{print ">"$2"\n"$3}' assembly/scinmit_hicul.hic.p_ctg.gfa > assembly/scinmit_hicul.hic.p_ctg.fa
```


Data inside `asm_scinmit/data` directory:
```bash
$ tree
.
├── hic
│   ├── lizard_hic_R1.fastq.gz -> /ibex/reference/course/cs249/lizard/input/hic/lizard_hic_R1.fastq.gz
│   └── lizard_hic_R2.fastq.gz -> /ibex/reference/course/cs249/lizard/input/hic/lizard_hic_R2.fastq.gz
├── ont
│   ├── lizard.fastq.gz -> /ibex/reference/course/cs249/lizard/input/ont/lizard.fastq.gz
│   ├── lizard_liver.fastq.gz -> /ibex/reference/course/cs249/lizard/input/ont/lizard_liver.fastq.gz
│   └── lizard_ont.fastq.gz -> /ibex/reference/course/cs249/lizard/input/ont/lizard_ont.fastq.gz
└── pacbio
    ├── lizard_liver.fastq.gz -> /ibex/reference/course/cs249/lizard/input/pacbio/lizard_liver.fastq.gz
    ├── lizard_liver_rev.fastq.gz -> /ibex/reference/course/cs249/lizard/input/pacbio/lizard_liver_rev.fastq.gz
    ├── lizard_liver_seq.fastq.gz -> /ibex/reference/course/cs249/lizard/input/pacbio/lizard_liver_seq.fastq.gz
    ├── lizard_rna_eye.fastq.gz -> /ibex/reference/course/cs249/lizard/input/pacbio/lizard_rna_eye.fastq.gz
    └── lizard_rna_liver.fastq.gz -> /ibex/reference/course/cs249/lizard/input/pacbio/lizard_rna_liver.fastq.gz
```

Primary contig assembly can be found in `/ibex/user/mohaah0a/cs249/BioAlgos-Assignment2/asm_scinmit/assembly/scinmit_hicul.hic.p_ctg.fa`.

### Task 2.2 Evaluation Scripts

#### QUAST
```bash 
# pip install quast
mkdir -p eval/quast
quast.py assembly/scinmit_hicul.hic.p_ctg.fa \
         --large \
         -t 64 \
         -o eval/quast
```

#### BUSCO
```bash
module load busco
# download the lineage once
busco --download sauropsida_odb10
# run Busco (took around 15 minutes)
busco -i assembly/scinmit_hicul.hic.p_ctg.fa \
      -l sauropsida_odb10 \
      -m genome \
      -o eval/busco \
      -c 64 \
      -f
```

#### Merqury
```bash
# first install meryl
wget https://github.com/marbl/meryl/releases/download/v1.4.1/meryl-1.4.1.Linux-amd64.tar.xz
tar -xJf meryl-1.4.1.Linux-amd64.tar.xz
export PATH=$PWD/meryl-1.4.1.Linux-amd64/bin:$PATH

mkdir -p logs/eval
# count read k-mers (took around 45 minutes)
meryl count k=21 threads=64 output eval/meryl/hifi_k21.meryl \
  data/pacbio/lizard_liver_seq.fastq.gz \
  data/pacbio/lizard_liver_rev.fastq.gz

# run Merqury (took around 17 minutes)
merqury.sh eval/meryl/hifi_k21.meryl assembly/scinmit_hicul.hic.p_ctg.fa eval/merqury
```

#### Inspector
```bash
# first install inspector
conda install -c bioconda inspector
git clone https://github.com/ChongLab/Inspector.git
export PATH=$PWD/Inspector/:$PATH

# run inspector (took around a day to finish)
inspector.py -c assembly/scinmit_hicul.hic.p_ctg.fa \
         -r data/pacbio/lizard_liver.fastq.gz data/ont/lizard_ont.fastq.gz \
         -o eval/inspector \
         -t 64 \
         --datatype mixed
```

## LLM usage

- I used Cursor IDE to write the code for this assignment.
  - Cursor has Cursor Tab functionality that auto completes lines and functions, which I used extensively but still manually reviewed every line written by it.
  - I used it while coding for all files.
  - It uses a mix of LLMs including GPT-4, Claude, Deepseek and Gemini.
- I used ChatGPT with the `gpt-4o` model to write:
  - test cases in earlier versions of `dbg_test.go` file.
  - `exportToGFA` function in `dbg.go` file.
  - `plot_kmer_spectra.py` script.
    - I reviewed the code outputs by running the scripts and comparing them with the expected outputs.
  - general questions related to the assignment and the bio-tools described in the assignment sheet.
  - Also, for help with running these bio-tools.
    - I reviewed the LLM outputs by running the commands and comparing them with the expected outputs.
- I used ChatGPT with the `o3` model to plan for the assembly of the lizard genome.
  - The initial plan has been heavily modified with each iteration of command running and output analysis.
- I tried using Aider with Claude 3.7 and gemini 2.5 flash to write functions in an early version of `dbg.go` file but all attempts failed.