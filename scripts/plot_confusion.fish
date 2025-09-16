#!/usr/bin/env fish

# Create confusion matrix plot directories
mkdir -p plots/normalized/confusion

# Function to create confusion matrix plot
function plot_confusion_matrix
    set group $argv[1]
    set data_file $argv[2]
    set output_file $argv[3]
    set title_text $argv[4]

    gnuplot -e "
    set terminal png size 700,500;
    set output '$output_file';
    set title '$title_text';
    set size ratio -1;
    set key noautotitle;
    set xrange [-0.5:1.5] noextend;
    set yrange [-0.5:1.5] reverse noextend;
    set xtics ('Big Car' 0, 'Small Car' 1);
    set ytics ('Big Car' 0, 'Small Car' 1);
    set xlabel 'Predicted';
    set ylabel 'Actual';
    set tics out;
    set palette rgb 33,13,10;
    set cbrange [0:*];
    set colorbox;
    plot '$data_file' using 2:1:3 with image, \
         '$data_file' using 2:1:(sprintf('%d',\$3)) with labels font 'Arial,14' tc 'black';
    "
end

# Plot confusion matrices for normalized data
if test -d plots/normalized/confusion
    for group in A B C
        set data_file "plots/normalized/confusion/group$group.dat"
        set output_file "plots/normalized/confusion/group$group.png"
        set title_text "Group $group - Normalized Data Confusion Matrix"

        if test -f $data_file
            plot_confusion_matrix $group $data_file $output_file $title_text
            echo "Generated: $output_file"
        else
            echo "Data file not found: $data_file"
        end
    end
end

echo "Confusion matrix plots generated:"
if test -d plots/normalized/confusion
    echo "  Normalized data: plots/normalized/confusion/"
    ls -la plots/normalized/confusion/ 2>/dev/null || echo "  No normalized confusion matrices found"
else
    echo "  No confusion matrix directory found"
end
