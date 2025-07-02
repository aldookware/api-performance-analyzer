#!/bin/sh
set -e

echo "🚀 Starting API Performance Analysis..."

# Get inputs from GitHub Action
CODE_PATH="${1:-${INPUT_CODE_PATH:-.}}"
OUTPUT_FORMAT="${2:-${INPUT_OUTPUT_FORMAT:-markdown}}"
SEVERITY_THRESHOLD="${3:-${INPUT_SEVERITY_THRESHOLD:-high}}"
FAIL_ON_ISSUES="${4:-${INPUT_FAIL_ON_ISSUES:-false}}"

echo "📁 Analyzing code in: $CODE_PATH"
echo "📊 Output format: $OUTPUT_FORMAT"
echo "🎯 Severity threshold: $SEVERITY_THRESHOLD"
echo "❌ Fail on issues: $FAIL_ON_ISSUES"
echo ""

# Set environment variables for the analyzer
export INPUT_CODE_PATH="$CODE_PATH"
export INPUT_OUTPUT_FORMAT="$OUTPUT_FORMAT"
export INPUT_SEVERITY_THRESHOLD="$SEVERITY_THRESHOLD"
export INPUT_FAIL_ON_ISSUES="$FAIL_ON_ISSUES"

# Run the analyzer
/root/analyzer \
  --path="$CODE_PATH" \
  --format="$OUTPUT_FORMAT" \
  --threshold="$SEVERITY_THRESHOLD" \
  --fail-on-issues="$FAIL_ON_ISSUES" \
  --verbose

EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
    echo ""
    echo "✅ Analysis complete! Check the results above."
else
    echo ""
    echo "❌ Analysis found issues above the specified threshold."
    echo "💡 Review the performance recommendations to improve your API."
fi

exit $EXIT_CODE
