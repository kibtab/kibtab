// K6 Load Test for Kibtab
// Usage: k6 run scripts/k6-load-test.js

export const options = {
  stages: [
    { duration: '30s', target: 10 },   // ramp up to 10 VUs
    { duration: '1m', target: 10 },    // stay at 10 VUs
    { duration: '30s', target: 0 },    // ramp down
  ],
  thresholds: {
    'http_req_duration': ['p(95)<500'], // 95% of requests under 500ms
    'http_req_failed': ['rate<0.01'],   // less than 1% failures
  },
};

export default function () {
  // Health check
  http.get('http://localhost:8080/healthz');

  // List tables
  http.get('http://localhost:8080/v1/tables');

  // Write a row
  http.post(
    'http://localhost:8080/v1/tables/my_table/rows',
    JSON.stringify({ row: 'load-test', cells: [{ field: 'k', value: 'k6-test' }] }),
    { headers: { 'Content-Type': 'application/json', 'X-Kibtab-User': 'loadtest' } }
  );
}