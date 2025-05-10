import matplotlib.pyplot as plt

filename = "merqury.scinmit_hicul.hic.p_ctg.spectra-cn.hist"
coverage = []
counts = []

with open(filename, 'r') as file:
    for line in file:
        if line.startswith('#') or line.lower().startswith('copies') or not line.strip():
            continue
        parts = line.strip().split()
        try:
            cov = int(parts[1])
            count = int(parts[2])
            coverage.append(cov)
            counts.append(count)
        except ValueError:
            continue

plt.figure(figsize=(12, 6))
plt.plot(coverage, counts, linewidth=1)
plt.title("K-mer Spectrum: Scincus mitranus")
plt.xlabel("K-mer Multiplicity (Coverage)")
plt.ylabel("Counts (Number of K-mers)")
plt.yscale("log")
plt.xlim(0, 1500)
plt.grid(True, linestyle='--', alpha=0.5)
plt.tight_layout()
plt.show()
