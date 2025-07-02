#!/bin/sh
set -e

# Get inputs from GitHub Action
CODE_PATH="${1:-${INPUT_CODE_PATH:-.}}"
OUTPUT_FORMAT="${2:-${INPUT_OUTPUT_FORMAT:-markdown}}"
SEVERITY_THRESHOLD="${3:-${INPUT_SEVERITY_THRESHOLD:-high}}"
FAIL_ON_ISSUES="${4:-${INPUT_FAIL_ON_ISSUES:-false}}"

# Only show verbose output for human-readable formats
if [ "$OUTPUT_FORMAT" = "markdown" ] || [ "$OUTPUT_FORMAT" = "github" ]; then
    echo "🚀 Starting API Performance Analysis..."
    echo "📁 Analyzing code in: $CODE_PATH"
    echo "📊 Output format: $OUTPUT_FORMAT"
    echo "🎯 Severity threshold: $SEVERITY_THRESHOLD"
    echo "❌ Fail on issues: $FAIL_ON_ISSUES"
    echo ""
    VERBOSE_FLAG="--verbose"
else
    # Silent mode for machine-readable formats (json, sarif)
    VERBOSE_FLAG=""
fi

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
  $VERBOSE_FLAG

EXIT_CODE=$?

# Only show completion messages for human-readable formats
if [ "$OUTPUT_FORMAT" = "markdown" ] || [ "$OUTPUT_FORMAT" = "github" ]; then
    if [ $EXIT_CODE -eq 0 ]; then
        echo ""
        echo "✅ Analysis complete! Check the results above."
    else
        echo ""
        echo "❌ Analysis found issues above the specified threshold."
        echo "💡 Review the performance recommendations to improve your API."
    fi
fi

exit $EXIT_CODE
