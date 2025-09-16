#!/usr/bin/env fish

echo "🚀 Starting Data Analysis Pipeline"
echo "=================================="

echo ""
echo "Step 1: Normalizing data..."
go run scripts/normalize.go
if test $status -ne 0
    echo "❌ Normalization failed"
    exit 1
end
echo "✅ Data normalization complete"

echo ""
echo "Step 2: Generating confusion matrices..."
go run scripts/confusion_matrix.go
if test $status -ne 0
    echo "❌ Confusion matrix analysis failed"
    exit 1
end
echo "✅ Confusion matrix analysis complete"

echo ""
echo "Step 3: Generating data plots..."
./scripts/plot.fish
if test $status -ne 0
    echo "❌ Data plotting failed"
    exit 1
end
echo "✅ Data plots generated"

echo ""
echo "Step 4: Generating confusion matrix plots..."
./scripts/plot_confusion.fish
if test $status -ne 0
    echo "❌ Confusion matrix plotting failed"
    exit 1
end
echo "✅ Confusion matrix plots generated"

echo ""
echo "Analysis pipeline complete!"
echo "Results available in:"
echo "  📁 plots/ - Original data plots"
echo "  📁 plots/normalized/ - Normalized data plots"
echo "  📁 plots/normalized/confusion/ - Confusion matrices and analysis results"
echo "  📄 plots/normalized/confusion/analysis_results.txt - Detailed metrics"
