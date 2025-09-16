# Data Processing & Visualization

Uses Nix for development environment.

## Setup
```bash
# With direnv
direnv allow

# Without direnv
nix develop
```

## File Structure
```
Data/
├── groupA.txt, groupB.txt, groupC.txt  # Raw CSV data
└── normalized/                         # Normalized output
scripts/
├── normalize.go                        # Data normalization
└── plot.fish                           # Plot generation
plots/                                  # Generated images
└── normalized/                         # Normalized data plots
flake.nix                               # Nix development environment
```

## Scripts
- `go run go/normalize.go` - Normalizes data to 0-1 range
- `fish scripts/plot.fish` - Generates plots with gnuplot