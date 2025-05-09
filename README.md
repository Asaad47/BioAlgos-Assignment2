# Assignment 2: Genome Assembly and Evaluation

- [Assignment 2: Genome Assembly and Evaluation](#assignment-2-genome-assembly-and-evaluation)
  - [Introduction](#introduction)
  - [LLM usage](#llm-usage)
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
      - [Flagger or Inspector](#flagger-or-inspector)
  - [Old scratch scripts](#old-scratch-scripts)
    - [Task 2.1](#task-21)
    - [Task 2.2](#task-22)


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

## Introduction

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
```bash
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_reads_hiseq_5k_dbg_k_40 reads/no_error_reads_hiseq_5k_dbg_k_40.fasta
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_ont_hq_50x_dbg_k_40 reads/no_error_ont_hq_50x_dbg_k_40.fasta

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_reads_hiseq_5k_olc_40 reads/no_error_reads_hiseq_5k_olc_40.fasta
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_ont_hq_50x_olc_40 reads/no_error_ont_hq_50x_olc_40.fasta
```

```bash
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_ont_hq_50x_dbg_k_40 reads/ont_hq_50x_dbg_k_40.fasta
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_reads_hiseq_5k_dbg_k_40 reads/reads_hiseq_5k_dbg_k_40.fasta

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_ont_hq_50x_olc_40 reads/ont_hq_50x_olc_40.fasta
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_reads_hiseq_5k_olc_40 reads/reads_hiseq_5k_olc_40.fasta
```

### Task 1.3.4 SPAdes
- Download SPAdes 4.1.0 (I couldn't get it from Ibex modules):
```bash
wget https://github.com/ablab/spades/releases/download/v4.1.0/SPAdes-4.1.0-Linux.tar.gz
 tar -xzf SPAdes-4.1.0-Linux.tar.gz
```

```bash
SPAdes-4.1.0-Linux/bin/spades.py --phred-offset 33 -s reads/no_error_reads_hiseq_5k.fastq -o spades_output_no_error_reads_hiseq_5k
SPAdes-4.1.0-Linux/bin/spades.py --phred-offset 33 -s reads/reads_hiseq_5k.fastq -o spades_output_hiseq_5k

SPAdes-4.1.0-Linux/bin/spades.py --phred-offset 33 -s reads/no_error_ont_hq_50x.fastq -o spades_output_no_error_ont_hq_50x
SPAdes-4.1.0-Linux/bin/spades.py --phred-offset 33 -s reads/ont_hq_50x.fastq -o spades_output_ont_hq_50x
```

```bash
quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_reads_hiseq_5k_spades spades_output_no_error_reads_hiseq_5k/contigs.fasta 

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_reads_hiseq_5k_spades spades_output_hiseq_5k/contigs.fasta 

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_no_error_ont_hq_50x_spades spades_output_no_error_ont_hq_50x/contigs.fasta 

quast.py -r GCF_000901155.1_ViralProj183710_genomic.fna -o quast_ont_hq_50x_spades spades_output_ont_hq_50x/contigs.fasta

```

### Task 2.1 Genome Assembly


Assmebly plan (based on ChatGPT o3 planning):
```bash
# 0- Allocate resources
# TODO: Add resources specification

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

#### Flagger or Inspector
```bash
# 4- Flagger or Inspector
```


inside `asm_scinmit/data` directory:
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
    ├── lizard_liver_rev_subsampled.fastq
    ├── lizard_liver_seq.fastq.gz -> /ibex/reference/course/cs249/lizard/input/pacbio/lizard_liver_seq.fastq.gz
    ├── lizard_liver_seq_subsampled.fastq
    ├── lizard_rna_eye.fastq.gz -> /ibex/reference/course/cs249/lizard/input/pacbio/lizard_rna_eye.fastq.gz
    └── lizard_rna_liver.fastq.gz -> /ibex/reference/course/cs249/lizard/input/pacbio/lizard_rna_liver.fastq.gz
```

TODOs:
- fix code for OLC on no-error and error reads.
- differentiate commands for setup and running the programs.


## Old scratch scripts

### Task 2.1

```bash
# 1- setup
module load gcc/12.2.0 hifiasm/0.19 minimap2 samtools yak meryl seqtk
# yak and meryl are not available in Ibex modules
mkdir -p asm_scinmit/{data,hifi,hic,ont,tmp,assembly}
cd asm_scinmit

# 2- fetch data
export RAW=/ibex/reference/course/cs249/lizard/input
# Create symbolic links inside project tree
ln -s $RAW/pacbio/*.fastq.gz data/pacbio/
ln -s $RAW/hic/*.fastq.gz data/hic/
ln -s $RAW/ont/*.fastq.gz data/ont/

# 4- optional downsampling (for testing)
seqtk sample -s42 data/pacbio/lizard_liver_seq.fastq.gz 0.10 > data/pacbio/lizard_liver_seq_subsampled.fastq # took 6 minutes
seqtk sample -s42 data/pacbio/lizard_liver_rev.fastq.gz 0.10 > data/pacbio/lizard_liver_rev_subsampled.fastq # took 8 minutes

# 5- primary assembly
# started running at 8:57 pm, finished 9:22 pm -> 25 minutes
hifiasm -o assembly/scinmit -t32 --primary -D10 data/pacbio/lizard_liver_seq_subsampled.fastq data/pacbio/lizard_liver_rev_subsampled.fastq

# mkdir -p full_assembly
# hifiasm -o full_assembly/scinmit -t32 --primary -D10 data/pacbio/lizard_liver_seq.fastq.gz data/pacbio/lizard_liver_rev.fastq.gz

# 6- convert GFA to FASTA
awk '$1=="S"{print ">"$2"\n"$3}' assembly/scinmit.p_ctg.gfa > assembly/scinmit.p_ctg.fa
awk '$1=="S"{print ">"$2"\n"$3}' assembly/scinmit.a_ctg.gfa > assembly/scinmit.a_ctg.fa



# 7- phase and scaffold with Hi-C
# Map Hi‑C reads
module load bwa
bwa index assembly/scinmit.p_ctg.fa # 23 minutes
bwa mem -5SP assembly/scinmit.p_ctg.fa data/hic/lizard_hic_R1.fastq.gz data/hic/lizard_hic_R2.fastq.gz > assembly/scinmit.hic.p_ctg.sam

# 8- (advanced) integrate ONT PromethION reads
hifiasm -o assembly/scinmit_ul -t32 --ul data/ont/lizard_ont.fastq.gz -l0 assembly/scinmit.hic.p_ctg.gfa



seqtk sample -s42 data/hic/lizard_hic_R1.fastq.gz 0.15 | gzip > tmp/hic_R1.15pct.fq.gz 
seqtk sample -s42 data/hic/lizard_hic_R2.fastq.gz 0.15 | gzip > tmp/hic_R2.15pct.fq.gz
export HIC_R1=tmp/hic_R1.15pct.fq.gz
export HIC_R2=tmp/hic_R2.15pct.fq.gz
#  ⬆ change to data/hic/lizard_hic_R[12].fastq.gz in the final run

# Run from inside asm_scinmit (with subsampled reads)
hifiasm -o assembly/scinmit_hicul_subsampled \
        -t32 --primary -D10 \
        --h1 $HIC_R1 \
        --h2 $HIC_R2 \
        --ul data/ont/lizard_ont.fastq.gz \
        data/pacbio/lizard_liver_seq_subsampled.fastq \
        data/pacbio/lizard_liver_rev_subsampled.fastq

# Run from inside asm_scinmit [ran at 11:29 pm Tue, finished 8:46 am Wed]
# Real time: 33467.576 sec; CPU: 1523075.897 sec; Peak RSS: 170.915 GB
hifiasm -o assembly/scinmit_hicul \
        -t64 --primary -D10 \
        --h1 data/hic/lizard_hic_R1.fastq.gz \
        --h2 data/hic/lizard_hic_R2.fastq.gz \
        --ul data/ont/lizard_ont.fastq.gz \
        data/pacbio/lizard_liver_seq.fastq.gz \
        data/pacbio/lizard_liver_rev.fastq.gz

awk '$1=="S"{print ">"$2"\n"$3}' assembly/scinmit_hicul.hic.p_ctg.gfa > assembly/scinmit_hicul.hic.p_ctg.fa & \
awk '$1=="S"{print ">"$2"\n"$3}' assembly/scinmit_hicul.hic.a_ctg.gfa > assembly/scinmit_hicul.hic.a_ctg.fa



# Ran at 12:17 pm Wed
hifiasm -o assembly/scinmit_hicul_dual.asm \
        -t64 --dual-scaf --primary \
        --h1 data/hic/lizard_hic_R1.fastq.gz \
        --h2 data/hic/lizard_hic_R2.fastq.gz \
        --ul data/ont/lizard_ont.fastq.gz \
        data/pacbio/lizard_liver_seq.fastq.gz \
        data/pacbio/lizard_liver_rev.fastq.gz \
        data/pacbio/lizard_liver.fastq.gz

# TODO: run the following
awk '$1=="S"{print ">"$2"\n"$3}' assembly/scinmit_hicul_dual.asm.hic.p_ctg.gfa > assembly/scinmit_hicul_dual.asm.hic.p_ctg.fa & \
awk '$1=="S"{print ">"$2"\n"$3}' assembly/scinmit_hicul_dual.asm.hic.a_ctg.gfa > assembly/scinmit_hicul_dual.asm.hic.a_ctg.fa

```


### Task 2.2

```bash
mkdir -p eval/{quast,busco,merqury,flagger}

# 1- Quast
module load quast/5.2.0
quast.py assembly/scinmit_hicul.hic.p_ctg.fa \
         --large \
         -t 64 \
         -o eval/quast

# 2- Busco
module load busco/5.4.7
# download the lineage once (∼1 GB) – replace sauropsida_odb10 if you know a closer clade
busco --download sauropsida_odb10
# started at 12:11 pm, finished 12:26 pm
busco -i assembly/scinmit_hicul.hic.p_ctg.fa \
      -l sauropsida_odb10 \
      -m genome \
      -o eval/busco \
      -c 64 \
      -f

# Output:
#  -------------------------------------------------------------------------------------------
#     |Results from dataset sauropsida_odb10                                                     |
#     -------------------------------------------------------------------------------------------
#     |C:98.2%[S:97.4%,D:0.9%],F:0.2%,M:1.5%,n:7480,E:2.4%                                       |
#     |7349    Complete BUSCOs (C)    (of which 178 contain internal stop codons)                |
#     |7285    Complete and single-copy BUSCOs (S)                                               |
#     |64    Complete and duplicated BUSCOs (D)                                                  |
#     |17    Fragmented BUSCOs (F)                                                               |
#     |114    Missing BUSCOs (M)                                                                 |
#     |7480    Total BUSCO groups searched                                                       |
#     -------------------------------------------------------------------------------------------
# 2025-05-08 12:26:23 INFO:       BUSCO analysis done with WARNING(s). Total running time: 882 seconds

# found in one_line_summary: "C:98.2%[S:97.4%,D:0.9%],F:0.2%,M:1.5%,n:7480,E:2.4%"

# one_line_summary: "C:98.2%[S:97.4%,D:0.9%],F:0.2%,M:1.5%,n:7480,E:2.4%"
# Complete percentage: 98.2
# Complete BUSCOs: 7349
# Single copy percentage: 97.4
# Single copy BUSCOs: 7285
# Multi copy percentage: 0.9
# Multi copy BUSCOs: 64
# Fragmented percentage: 0.2
# Fragmented BUSCOs: 17
# Missing percentage: 1.5
# Missing BUSCOs: 114
# n_markers: 7480
# avg_identity: 0.81
# domain: "eukaryota"
# internal_stop_codon_count: 178
# internal_stop_codon_percent: 2.4

# 3- Merqury

wget https://github.com/marbl/meryl/releases/download/v1.4.1/meryl-1.4.1.Linux-amd64.tar.xz
tar -xJf meryl-1.4.1.Linux-amd64.tar.xz
cd meryl-1.4.1.Linux-amd64/bin
export PATH=$PWD:$PATH

# module load merqury/1.3 meryl/1.3
# # 3.1  count read k‑mers (takes ~2 h, 250 GB RAM for 30× HiFi)
meryl count k=21 output eval/meryl/hifi_k21_hicul_p.meryl \
      assembly/scinmit_hicul.hic.p_ctg.fa

# # 3.2  run Merqury
rm -r eval/merqury
mkdir -p logs/eval
# ran at 11:07 am Wed, finished 11:15 am
merqury.sh eval/meryl/hifi_k21_hicul_p.meryl assembly/scinmit_hicul.hic.p_ctg.fa eval/merqury
# gave errored outputs


# started at 11:58 am, finished at 12:41 **pm**
meryl count k=21 threads=64 output eval/meryl/hifi_k21.meryl \
  data/pacbio/lizard_liver_seq.fastq.gz \
  data/pacbio/lizard_liver_rev.fastq.gz

# started at 2:22 pm, finished at 2:39 pm
merqury.sh eval/meryl/hifi_k21.meryl assembly/scinmit_hicul.hic.p_ctg.fa eval/merqury
# QV 65.6166 ????



git clone https://github.com/lh3/yak.git
cd yak && make -j && export PATH=$PWD:$PATH

mkdir -p eval/yak
yak count -k21 -b37 -t32 -o eval/yak/hifi_k21.yak data/pacbio/lizard_liver_seq_subsampled.fastq data/pacbio/lizard_liver_rev_subsampled.fastq
# Real time: 127.475 sec; CPU: 1181.206 sec; Peak RSS: 24.428 GB

yak qv -t32 -p eval/yak/hifi_k21.yak assembly/scinmit.p_ctg.fa > eval/yak/qv.txt
# Real time: 27.962 sec; CPU: 231.020 sec; Peak RSS: 9.658 GB
# QV 16.306

yak count -k21 -b37 -t64 -o eval/yak/hifi_k21_hicul.yak assembly/scinmit_hicul.hic.p_ctg.fa assembly/scinmit_hicul.hic.a_ctg.fa
# Real time: 45.984 sec; CPU: 404.851 sec; Peak RSS: 21.149 GB
yak qv -t64 -p eval/yak/hifi_k21_hicul.yak assembly/scinmit_hicul.hic.p_ctg.fa > eval/yak/qv_hicul.txt
# Real time: 105.463 sec; CPU: 281.682 sec; Peak RSS: 10.542 GB
# QV 30.950

yak count -k21 -b37 -t64 -o eval/yak/hifi_k21_hicul_p.yak assembly/scinmit_hicul.hic.p_ctg.fa
# Real time: 53.367 sec; CPU: 406.691 sec; Peak RSS: 21.139 GB
yak qv -t64 -p eval/yak/hifi_k21_hicul_p.yak assembly/scinmit_hicul.hic.p_ctg.fa > eval/yak/qv_hicul_p.txt
# Real time: 86.505 sec; CPU: 236.300 sec; Peak RSS: 11.058 GB
# QV 34.370

yak count -k30 -b37 -t64 -o eval/yak/hifi_k30_hicul_p.yak assembly/scinmit_hicul.hic.p_ctg.fa
# Real time: 52.007 sec; CPU: 422.721 sec; Peak RSS: 20.902 GB
yak qv -t64 -p eval/yak/hifi_k30_hicul_p.yak assembly/scinmit_hicul.hic.p_ctg.fa > eval/yak/qv_hicul_p_k30.txt
# Real time: 75.037 sec; CPU: 214.332 sec; Peak RSS: 10.542 GB
# QV 34.165

yak count -k25 -b37 -t64 -o eval/yak/hifi_k25_hicul_p.yak assembly/scinmit_hicul.hic.p_ctg.fa
# Real time: 53.093 sec; CPU: 430.635 sec; Peak RSS: 21.087 GB
yak qv -t64 -p eval/yak/hifi_k25_hicul_p.yak assembly/scinmit_hicul.hic.p_ctg.fa > eval/yak/qv_hicul_p_k25.txt
# Real time: 87.560 sec; CPU: 248.112 sec; Peak RSS: 10.545 GB
# QV 34.277


# count should be done on the reads not the assembly
yak count -k21 -b37 -t64 -o eval/yak/hifi_k21_hicul_p_reads.yak data/pacbio/lizard_liver_seq.fastq.gz data/pacbio/lizard_liver_rev.fastq.gz
# Real time: 1390.887 sec; CPU: 33102.092 sec; Peak RSS: 48.797 GB
yak qv -t64 -p eval/yak/hifi_k21_hicul_p_reads.yak assembly/scinmit_hicul.hic.p_ctg.fa > eval/yak/qv_hicul_p_reads.txt
# Real time: 192.922 sec; CPU: 529.968 sec; Peak RSS: 42.528 GB
# QV 40.436

# using same base count for reads
yak qv -t64 -p eval/yak/hifi_k21_hicul_p_reads.yak assembly/scinmit_hicul_dual.asm.hic.p_ctg.fa > eval/yak/qv_hicul_p_reads_dual.txt
# Real time: 204.390 sec; CPU: 834.617 sec; Peak RSS: 50.365 GB
# QV 40.376

yak count -k21 -b37 -t64 -o eval/yak/hifi_k21_hicul_p_reads_dual.yak data/pacbio/lizard_liver_seq.fastq.gz data/pacbio/lizard_liver_rev.fastq.gz data/pacbio/lizard_liver.fastq.gz
# Real time: 1086.267 sec; CPU: 17299.753 sec; Peak RSS: 48.792 GB
yak qv -t64 -p eval/yak/hifi_k21_hicul_p_reads_dual.yak assembly/scinmit_hicul_dual.asm.hic.p_ctg.fa > eval/yak/qv_hicul_p_reads_dual_extended.txt
# Real time: 218.194 sec; CPU: 967.875 sec; Peak RSS: 50.165 GB
# QV 40.376




# 4- Flagger
# module load winnowmap/2.1 flagger/0.2 samtools/1.17
module load winnowmap flagger samtools
# 4.1  build a Winnowmap repeat‑mask (speeds up mapping)
meryl count k=15 output repeats.meryl assembly/scinmit.p_ctg.fa
meryl print greater-than 1000 repeats.meryl > repeats_k15.txt

# 4.2  map HiFi reads
winnowmap -W repeats_k15.txt -ax map-pb -t 32 \
          assembly/scinmit.p_ctg.fa \
          /ibex/reference/course/lizard/input/pacbio/hifi/*.fastq.gz | \
    samtools sort -@32 -o eval/flagger/hifi.bam
samtools index eval/flagger/hifi.bam

# 4.3  run Flagger
flagger --bam eval/flagger/hifi.bam \
        --asm assembly/scinmit.p_ctg.fa \
        --outdir eval/flagger \
        --threads 32

# 4.alternative- Inspector
module load minimap2 samtools 
minimap2 -d scinmit.mmi assembly/scinmit_hicul.hic.p_ctg.fa
# Real time: 38.265 sec; CPU: 63.051 sec; Peak RSS: 7.181 GB

minimap2 -ax map-hifi scinmit.mmi data/pacbio/lizard_liver.fastq.gz | samtools sort -o aln.sorted.bam
samtools index aln.sorted.bam

inspector -r assembly/scinmit_hicul.hic.p_ctg.fa \
          -b aln.sorted.bam \
          -o inspector_out \
          -t 64


# started at 1:31 pm
inspector.py -c assembly/scinmit_hicul.hic.p_ctg.fa \
         -r data/pacbio/lizard_liver.fastq.gz data/ont/lizard_ont.fastq.gz \
         -o eval/inspector \
         -t 64 \
         --datatype mixed
# partial: Real time: 2570.461 sec; CPU: 137914.407 sec; Peak RSS: 20.368 GB
# partial: Real time: 2450.616 sec; CPU: 67002.273 sec; Peak RSS: 26.257 GB

```