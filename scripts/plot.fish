#!/usr/bin/env fish

# Create plots directory
mkdir -p plots plots/normalized

# Plot original data - all groups
gnuplot -e "
set terminal png size 800,600;
set output 'plots/original_data.png';
set title 'Original Data - All Groups';
set xlabel 'X';
set ylabel 'Y';
set key outside;
set datafile separator ',';
plot 'Data/groupA.txt' using 1:2:3 with points pt 7 ps 0.5 lc variable title 'Group A', \
     'Data/groupB.txt' using 1:2:3 with points pt 9 ps 0.5 lc variable title 'Group B', \
     'Data/groupC.txt' using 1:2:3 with points pt 11 ps 0.5 lc variable title 'Group C';
"

# Individual original group plots
for group in A B C
    gnuplot -e "
    set terminal png size 600,600;
    set output 'plots/group$group.png';
    set title 'Group $group Data Points';
    set xlabel 'X';
    set ylabel 'Y';
    set datafile separator ',';
    plot 'Data/group$group.txt' using 1:2:3 with points pt 7 ps 0.8 lc variable notitle;
    "
end

# Check if normalized data exists
if test -d Data/normalized
    # Plot normalized data - all groups
    gnuplot -e "
    set terminal png size 800,600;
    set output 'plots/normalized/normalized_data.png';
    set title 'Normalized Data (0-1 range) - All Groups';
    set xlabel 'X (normalized)';
    set ylabel 'Y (normalized)';
    set key outside;
    set datafile separator ',';
    plot 'Data/normalized/groupA.txt' using 1:2:3 with points pt 7 ps 0.5 lc variable title 'Group A', \
         'Data/normalized/groupB.txt' using 1:2:3 with points pt 9 ps 0.5 lc variable title 'Group B', \
         'Data/normalized/groupC.txt' using 1:2:3 with points pt 11 ps 0.5 lc variable title 'Group C';
    "

    # Individual normalized group plots
    for group in A B C
        gnuplot -e "
        set terminal png size 600,600;
        set output 'plots/normalized/group$group.png';
        set title 'Group $group Normalized Data Points';
        set xlabel 'X (normalized)';
        set ylabel 'Y (normalized)';
        set datafile separator ',';
        plot 'Data/normalized/group$group.txt' using 1:2:3 with points pt 7 ps 0.8 lc variable notitle;
        "
    end
end

echo "Plots generated:"
echo "  Original data: plots/"
if test -d Data/normalized
    echo "  Normalized data: plots/normalized/"
end
ls -la plots/
if test -d plots/normalized
    ls -la plots/normalized/
end
