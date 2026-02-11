/**
 * CPS Action API Performance Test Script
 * Parses Postman collection and tests every endpoint against localhost
 * 
 * Usage: node scripts/api_performance_test.js [--token YOUR_JWT_TOKEN]
 */

const fs = require('fs');
const http = require('http');
const path = require('path');

// ─── Configuration ───────────────────────────────────────────────────────────
const BASE_URL = 'http://localhost:8080/api/v1/cbesuperapp';
const POSTMAN_FILE = path.join(__dirname, '..', 'NEW FULL CPS ACTION.postman_collection.json');
const TIMEOUT_MS = 3000;
const CONCURRENCY = 1; // Sequential to avoid overwhelming the server

// Parse --token argument or read from .auth_token file
let AUTH_TOKEN = '';
const tokenIdx = process.argv.indexOf('--token');
if (tokenIdx !== -1 && process.argv[tokenIdx + 1]) {
  AUTH_TOKEN = process.argv[tokenIdx + 1];
} else {
  const tokenFile = path.join(__dirname, '.auth_token');
  if (fs.existsSync(tokenFile)) {
    AUTH_TOKEN = fs.readFileSync(tokenFile, 'utf8').trim();
  }
}

// ─── Variable Mapping ────────────────────────────────────────────────────────
// All Postman env variables map to the same local base
const VARIABLE_MAP = {
  '{{LOCAL}}': BASE_URL,
  '{{local}}': BASE_URL,
  '{{DEV}}': BASE_URL,
  '{{dev}}': BASE_URL,
  '{{dev_test}}': BASE_URL,
  '{{QA}}': BASE_URL,
  '{{qa}}': BASE_URL,
  '{{core_qa}}': BASE_URL,
  '{{CORE-QA}}': BASE_URL,
  '{{COREQA}}': BASE_URL,
  '{{VM-LOCAL}}': BASE_URL,
  // Placeholder IDs for parameterized routes
  '{{id}}': 'test-id-123',
  '{{ID}}': 'test-id-123',
  '{{bankID}}': 'test-bank-id',
  '{{branch_code}}': 'BR001',
  '{{region_code}}': 'AA',
  '{{ObjectId}}': '000000000000000000000001',
  '{{METHOD}}': 'test-method',
  '{{action_code}}': 'test-action-code',
  '{{user_code}}': 'test-user-code',
  '{{user_id}}': 'test-user-id',
  '{{customer_number}}': '12345',
  '{{service_code}}': 'test-service-code',
  '{{code}}': 'test-code',
  '{{merchant_id}}': 'test-merchant-id',
  '{{cps-maker-token}}': AUTH_TOKEN,
  '{{cps-checker-token}}': AUTH_TOKEN,
};

// ─── Postman Collection Parser ───────────────────────────────────────────────
function parseCollection(filePath) {
  const raw = fs.readFileSync(filePath, 'utf8');
  const collection = JSON.parse(raw);
  const endpoints = [];

  function walk(items, folderPath) {
    for (const item of items) {
      if (item.request) {
        const req = item.request;
        const method = req.method || 'GET';
        let rawUrl = req.url?.raw || '';
        
        // Skip empty URLs
        if (!rawUrl) continue;

        // Resolve variables
        let resolvedUrl = rawUrl;
        for (const [key, val] of Object.entries(VARIABLE_MAP)) {
          resolvedUrl = resolvedUrl.split(key).join(val);
        }
        // Handle any remaining unresolved {{vars}}
        resolvedUrl = resolvedUrl.replace(/\{\{[^}]+\}\}/g, 'test-placeholder');

        // Get body
        let body = '';
        if (req.body && req.body.raw) {
          body = req.body.raw;
          // Remove JS-style comments from JSON body
          body = body.replace(/\/\/.*$/gm, '').trim();
        }

        // Check if it's a multipart form
        const isMultipart = req.body?.mode === 'formdata';

        endpoints.push({
          name: item.name,
          method,
          rawUrl,
          resolvedUrl,
          body,
          isMultipart,
          folder: folderPath,
          hasAuth: !!(req.auth || AUTH_TOKEN),
        });
      }
      if (item.item) {
        walk(item.item, folderPath + ' > ' + item.name);
      }
    }
  }

  walk(collection.item, '');
  return endpoints;
}

// ─── HTTP Request with Timing ────────────────────────────────────────────────
function makeRequest(method, urlStr, body, token) {
  return new Promise((resolve) => {
    const startTime = process.hrtime.bigint();
    
    let parsedUrl;
    try {
      parsedUrl = new URL(urlStr);
    } catch (e) {
      resolve({
        status: 0,
        statusText: 'INVALID_URL',
        responseTime: 0,
        error: `Invalid URL: ${urlStr}`,
        bodyPreview: '',
      });
      return;
    }

    const options = {
      hostname: parsedUrl.hostname,
      port: parsedUrl.port || 80,
      path: parsedUrl.pathname + parsedUrl.search,
      method: method,
      timeout: TIMEOUT_MS,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    };

    if (token) {
      options.headers['Authorization'] = `Bearer ${token}`;
    }

    if (body && ['POST', 'PUT', 'PATCH'].includes(method)) {
      options.headers['Content-Length'] = Buffer.byteLength(body);
    }

    const req = http.request(options, (res) => {
      let data = '';
      res.on('data', (chunk) => { data += chunk; });
      res.on('end', () => {
        const endTime = process.hrtime.bigint();
        const responseTime = Number(endTime - startTime) / 1e6; // ms
        resolve({
          status: res.statusCode,
          statusText: res.statusMessage,
          responseTime: Math.round(responseTime * 100) / 100,
          error: null,
          bodyPreview: data.substring(0, 200),
          bodySize: data.length,
        });
      });
    });

    req.on('error', (err) => {
      const endTime = process.hrtime.bigint();
      const responseTime = Number(endTime - startTime) / 1e6;
      resolve({
        status: 0,
        statusText: 'ERROR',
        responseTime: Math.round(responseTime * 100) / 100,
        error: err.message,
        bodyPreview: '',
      });
    });

    req.on('timeout', () => {
      req.destroy();
      const endTime = process.hrtime.bigint();
      const responseTime = Number(endTime - startTime) / 1e6;
      resolve({
        status: 0,
        statusText: 'TIMEOUT',
        responseTime: Math.round(responseTime * 100) / 100,
        error: 'Request timed out',
        bodyPreview: '',
      });
    });

    if (body && ['POST', 'PUT', 'PATCH'].includes(method)) {
      req.write(body);
    }
    req.end();
  });
}

// ─── Results Formatter ───────────────────────────────────────────────────────
function formatResults(results) {
  const summary = {
    total: results.length,
    success: 0,    // 2xx
    authError: 0,  // 401/403
    clientError: 0, // 4xx (non-auth)
    serverError: 0, // 5xx
    networkError: 0, // 0 (connection refused, timeout)
    avgResponseTime: 0,
    minResponseTime: Infinity,
    maxResponseTime: 0,
    totalTime: 0,
  };

  for (const r of results) {
    if (r.result.status >= 200 && r.result.status < 300) summary.success++;
    else if (r.result.status === 401 || r.result.status === 403) summary.authError++;
    else if (r.result.status >= 400 && r.result.status < 500) summary.clientError++;
    else if (r.result.status >= 500) summary.serverError++;
    else summary.networkError++;

    if (r.result.responseTime > 0) {
      summary.totalTime += r.result.responseTime;
      if (r.result.responseTime < summary.minResponseTime) summary.minResponseTime = r.result.responseTime;
      if (r.result.responseTime > summary.maxResponseTime) summary.maxResponseTime = r.result.responseTime;
    }
  }
  summary.avgResponseTime = Math.round((summary.totalTime / results.length) * 100) / 100;

  return summary;
}

function getStatusEmoji(status) {
  if (status >= 200 && status < 300) return 'OK';
  if (status === 401 || status === 403) return 'AUTH';
  if (status >= 400 && status < 500) return 'CLI_ERR';
  if (status >= 500) return 'SRV_ERR';
  return 'NET_ERR';
}

function getPerformanceRating(ms) {
  if (ms <= 100) return 'FAST';
  if (ms <= 500) return 'GOOD';
  if (ms <= 1000) return 'SLOW';
  return 'VERY_SLOW';
}

// ─── Main ────────────────────────────────────────────────────────────────────
async function main() {
  console.log('='.repeat(100));
  console.log('  CPS ACTION API - PERFORMANCE TEST');
  console.log('  Base URL:', BASE_URL);
  console.log('  Auth Token:', AUTH_TOKEN ? 'PROVIDED (' + AUTH_TOKEN.substring(0, 20) + '...)' : 'NOT PROVIDED (will get 401 on protected routes)');
  console.log('='.repeat(100));

  // Step 1: Check server is reachable
  console.log('\n[1/4] Checking server connectivity...');
  const healthCheck = await makeRequest('GET', BASE_URL + '/cps_action/healthcheck', '', '');
  if (healthCheck.status === 0) {
    console.error('  ERROR: Server is not reachable at', BASE_URL);
    console.error('  Error:', healthCheck.error);
    console.error('  Please start the server first: go run cmd/main.go');
    process.exit(1);
  }
  console.log('  Server is UP! Health check returned:', healthCheck.status, '-', healthCheck.bodyPreview.substring(0, 100));

  // Step 2: Parse Postman collection
  console.log('\n[2/4] Parsing Postman collection...');
  const endpoints = parseCollection(POSTMAN_FILE);
  console.log(`  Found ${endpoints.length} endpoints`);

  // Deduplicate by resolved URL + method (some endpoints are duplicated across folders)
  const seen = new Set();
  const uniqueEndpoints = endpoints.filter(ep => {
    const key = ep.method + ' ' + ep.resolvedUrl;
    if (seen.has(key)) return false;
    seen.add(key);
    return true;
  });
  console.log(`  Unique endpoints (after dedup): ${uniqueEndpoints.length}`);

  // Filter out non-cps_action endpoints (like cps_auth ones)
  const cpsActionEndpoints = uniqueEndpoints.filter(ep => ep.resolvedUrl.includes('/cps_action/'));
  const authEndpoints = uniqueEndpoints.filter(ep => ep.resolvedUrl.includes('/cps_auth/') || ep.resolvedUrl.includes('/auth/'));
  console.log(`  CPS Action endpoints: ${cpsActionEndpoints.length}`);
  console.log(`  Auth endpoints (skipped): ${authEndpoints.length}`);

  // Step 3: Test all endpoints
  console.log('\n[3/4] Testing all CPS Action endpoints...\n');
  console.log(padRight('#', 4) + padRight('METHOD', 8) + padRight('STATUS', 8) + padRight('TIME(ms)', 10) + padRight('PERF', 12) + padRight('CATEGORY', 10) + padRight('ENDPOINT NAME', 45) + 'URL PATH');
  console.log('-'.repeat(150));

  const results = [];
  let idx = 0;

  for (const ep of cpsActionEndpoints) {
    idx++;
    // Skip multipart/form-data requests (file uploads can't be tested easily)
    if (ep.isMultipart) {
      results.push({
        endpoint: ep,
        result: { status: 0, statusText: 'SKIPPED', responseTime: 0, error: 'Multipart form (file upload)', bodyPreview: '' }
      });
      console.log(padRight(idx.toString(), 4) + padRight(ep.method, 8) + padRight('SKIP', 8) + padRight('-', 10) + padRight('-', 12) + padRight('SKIP', 10) + padRight(ep.name.substring(0, 43), 45) + getShortPath(ep.resolvedUrl));
      continue;
    }

    const result = await makeRequest(ep.method, ep.resolvedUrl, ep.body, AUTH_TOKEN);
    results.push({ endpoint: ep, result });

    const statusStr = result.status > 0 ? result.status.toString() : 'ERR';
    const timeStr = result.responseTime > 0 ? result.responseTime.toString() : '-';
    const category = getStatusEmoji(result.status);
    const perf = result.responseTime > 0 ? getPerformanceRating(result.responseTime) : '-';

    console.log(
      padRight(idx.toString(), 4) +
      padRight(ep.method, 8) +
      padRight(statusStr, 8) +
      padRight(timeStr, 10) +
      padRight(perf, 12) +
      padRight(category, 10) +
      padRight(ep.name.substring(0, 43), 45) +
      getShortPath(ep.resolvedUrl)
    );
  }

  // Step 4: Summary
  console.log('\n' + '='.repeat(100));
  console.log('  TEST RESULTS SUMMARY');
  console.log('='.repeat(100));

  const summary = formatResults(results);
  console.log(`\n  Total Endpoints Tested:  ${summary.total}`);
  console.log(`  ────────────────────────────────────`);
  console.log(`  2xx Success:             ${summary.success}  (${pct(summary.success, summary.total)})`);
  console.log(`  401/403 Auth Required:   ${summary.authError}  (${pct(summary.authError, summary.total)})`);
  console.log(`  4xx Client Errors:       ${summary.clientError}  (${pct(summary.clientError, summary.total)})`);
  console.log(`  5xx Server Errors:       ${summary.serverError}  (${pct(summary.serverError, summary.total)})`);
  console.log(`  Network/Timeout Errors:  ${summary.networkError}  (${pct(summary.networkError, summary.total)})`);
  console.log(`\n  PERFORMANCE METRICS`);
  console.log(`  ────────────────────────────────────`);
  console.log(`  Average Response Time:   ${summary.avgResponseTime} ms`);
  console.log(`  Min Response Time:       ${summary.minResponseTime === Infinity ? 'N/A' : summary.minResponseTime} ms`);
  console.log(`  Max Response Time:       ${summary.maxResponseTime} ms`);
  console.log(`  Total Test Duration:     ${Math.round(summary.totalTime)} ms`);

  // Performance distribution
  const fast = results.filter(r => r.result.responseTime > 0 && r.result.responseTime <= 100).length;
  const good = results.filter(r => r.result.responseTime > 100 && r.result.responseTime <= 500).length;
  const slow = results.filter(r => r.result.responseTime > 500 && r.result.responseTime <= 1000).length;
  const verySlow = results.filter(r => r.result.responseTime > 1000).length;
  console.log(`\n  PERFORMANCE DISTRIBUTION`);
  console.log(`  ────────────────────────────────────`);
  console.log(`  FAST    (<=100ms):   ${fast} endpoints`);
  console.log(`  GOOD    (<=500ms):   ${good} endpoints`);
  console.log(`  SLOW    (<=1000ms):  ${slow} endpoints`);
  console.log(`  VERY_SLOW (>1000ms): ${verySlow} endpoints`);

  // List server errors (5xx)
  const serverErrors = results.filter(r => r.result.status >= 500);
  if (serverErrors.length > 0) {
    console.log(`\n  SERVER ERRORS (5xx) - NEEDS ATTENTION`);
    console.log(`  ────────────────────────────────────`);
    for (const r of serverErrors) {
      console.log(`  [${r.result.status}] ${r.endpoint.method} ${r.endpoint.name}`);
      console.log(`       URL: ${getShortPath(r.endpoint.resolvedUrl)}`);
      console.log(`       Response: ${r.result.bodyPreview.substring(0, 150)}`);
    }
  }

  // List slow endpoints
  const slowEndpoints = results.filter(r => r.result.responseTime > 500).sort((a, b) => b.result.responseTime - a.result.responseTime);
  if (slowEndpoints.length > 0) {
    console.log(`\n  SLOW ENDPOINTS (>500ms) - PERFORMANCE CONCERN`);
    console.log(`  ────────────────────────────────────`);
    for (const r of slowEndpoints.slice(0, 20)) {
      console.log(`  [${r.result.responseTime}ms] ${r.endpoint.method} ${r.endpoint.name} - ${getShortPath(r.endpoint.resolvedUrl)}`);
    }
  }

  // Write detailed results to JSON
  const reportPath = path.join(__dirname, '..', 'api_test_report.json');
  const report = {
    timestamp: new Date().toISOString(),
    baseUrl: BASE_URL,
    authProvided: !!AUTH_TOKEN,
    summary,
    endpoints: results.map(r => ({
      name: r.endpoint.name,
      method: r.endpoint.method,
      url: getShortPath(r.endpoint.resolvedUrl),
      folder: r.endpoint.folder,
      status: r.result.status,
      responseTime: r.result.responseTime,
      category: getStatusEmoji(r.result.status),
      performance: r.result.responseTime > 0 ? getPerformanceRating(r.result.responseTime) : 'N/A',
      error: r.result.error,
      responsePreview: r.result.bodyPreview?.substring(0, 200),
    })),
  };
  fs.writeFileSync(reportPath, JSON.stringify(report, null, 2));
  console.log(`\n  Full report saved to: ${reportPath}`);
  console.log('='.repeat(100));
}

// ─── Helpers ─────────────────────────────────────────────────────────────────
function padRight(str, len) {
  str = String(str);
  return str.length >= len ? str.substring(0, len) : str + ' '.repeat(len - str.length);
}

function pct(val, total) {
  return total > 0 ? Math.round((val / total) * 100) + '%' : '0%';
}

function getShortPath(url) {
  try {
    const u = new URL(url);
    return u.pathname + u.search;
  } catch {
    return url;
  }
}

// ─── Run ─────────────────────────────────────────────────────────────────────
main().catch(err => {
  console.error('Fatal error:', err);
  process.exit(1);
});
