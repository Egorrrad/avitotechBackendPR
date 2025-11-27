#!/bin/bash

echo "=== Load Testing ==="
echo ""

if ! command -v k6 &> /dev/null; then
    echo "k6 не установлен"
    exit 1
fi

BASE_URL=${BASE_URL:-"http://localhost:8080"}
OUTPUT_DIR=${OUTPUT_DIR:-"tests/load/results"}

echo "${BASE_URL}"
curl -f "${BASE_URL}/health" || {
    echo "Service is not healthy"
    exit 1
}

mkdir -p "$OUTPUT_DIR"

echo "Запуск теста..."
K6_WEB_DASHBOARD=true K6_WEB_DASHBOARD_EXPORT="${OUTPUT_DIR}/report.html" \
k6 run \
  --out json="${OUTPUT_DIR}/results.json" \
  --summary-export="${OUTPUT_DIR}/summary.json" \
  -e BASE_URL="$BASE_URL" \
  "$(dirname "$0")/load-test.js"

EXIT_CODE=$?

echo ""
if [ $EXIT_CODE -eq 0 ]; then
    echo "   Тест завершён успешно!"
else
    echo "   Тест не прошёл!"
fi

echo ""
echo "Результаты сохранены в: $OUTPUT_DIR"
echo ""

exit $EXIT_CODE