import http from 'k6/http';
import {check} from 'k6';
import {Rate, Trend} from 'k6/metrics';

const errorRate = new Rate('errors');
const responseTime = new Trend('response_time');

export const options = {
    stages: [
        { duration: '30s', target: 5 },   // Разогрев
        { duration: '5m', target: 5 },    // Стабильная нагрузка 5 RPS
        { duration: '30s', target: 0 },   // Плавное завершение
    ],
    thresholds: {
        'http_req_duration': ['p(95)<300'],
        'errors': ['rate<0.001'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
    console.log('Setting up test data: 20 teams, 200 users...');
    const teams = [];
    const users = [];

    for (let i = 1; i <= 20; i++) {
        const teamName = `team_${crypto.randomUUID()}`;
        const members = [];

        for (let j = 1; j <= 10; j++) {
            const userId = `user_${crypto.randomUUID()}`;
            members.push({
                user_id: userId,
                username: `User_${i}_${j}`,
                is_active: true
            });
            users.push({ user_id: userId, team_name: teamName });
        }

        const response = http.post(
            `${BASE_URL}/team/add`,
            JSON.stringify({
                team_name: teamName,
                members: members
            }),
            {
                headers: { 'Content-Type': 'application/json' },
                tags: { name: 'setup' }
            }
        );

        check(response, { 'setup team: status 201': (r) => r.status === 201 });
        if (response.status === 201) {
            teams.push(teamName);
        }
    }

    console.log(`Setup complete: ${teams.length} teams, ${users.length} users`);
    return { teams, users };
}

export default function(data) {
    const { teams, users } = data;

    const randomUser = users[Math.floor(Math.random() * users.length)];
    const randomTeam = teams[Math.floor(Math.random() * teams.length)];
    const prId = `pr_${crypto.randomUUID()}`;

    const rand = Math.random();

    // 5% — Получение команды (GET /team/get)
    if (rand < 0.05) {
        if (teams.length === 0) return;

        const res = http.get(
            `${BASE_URL}/team/get?team_name=${encodeURIComponent(randomTeam)}`,
            { tags: { name: 'get_team' } }
        );

        const success = check(res, {
            'get_team: status 200': (r) => r.status === 200,
            'get_team: valid JSON with team_name and members': (r) => {
                try {
                    const body = JSON.parse(r.body);
                    return body.team && body.team.team_name === teamName && Array.isArray(body.team.members);
                } catch (e) {
                    return false;
                }
            },
            'get_team: response time < 300ms': (r) => r.timings.duration < 300,
        });

        errorRate.add(!success);
        responseTime.add(res.timings.duration, { operation: 'get_team' });
    }

    // 35% — Создание PR (POST /pullRequest/create)
    else if (rand < 0.40) {
        const res = http.post(
            `${BASE_URL}/pullRequest/create`,
            JSON.stringify({
                pull_request_id: prId,
                pull_request_name: `Feature ${prId}`,
                author_id: randomUser.user_id
            }),
            {
                headers: { 'Content-Type': 'application/json' },
                tags: { name: 'create_pr' }
            }
        );

        const success = check(res, {
            'create_pr: status 201 or 409': (r) => r.status === 201 || r.status === 409,
            'create_pr: response time < 300ms': (r) => r.timings.duration < 300,
        });

        errorRate.add(!success);
        responseTime.add(res.timings.duration, { operation: 'create_pr' });
    }

    // 25% — Получение списка PR для ревью (GET /users/getReview)
    else if (rand < 0.65) {
        const res = http.get(
            `${BASE_URL}/users/getReview?user_id=${randomUser.user_id}`,
            { tags: { name: 'get_review' } }
        );

        const success = check(res, {
            'get_review: status 200': (r) => r.status === 200,
            'get_review: valid JSON': (r) => {
                try {
                    const body = JSON.parse(r.body);
                    return body.hasOwnProperty('user_id') && body.hasOwnProperty('pull_requests');
                } catch (e) {
                    return false;
                }
            },
            'get_review: response time < 300ms': (r) => r.timings.duration < 300,
        });

        errorRate.add(!success);
        responseTime.add(res.timings.duration, { operation: 'get_review' });
    }

    // 15% — Merge PR (POST /pullRequest/merge)
    else if (rand < 0.80) {
        const res = http.post(
            `${BASE_URL}/pullRequest/merge`,
            JSON.stringify({ pull_request_id: prId }),
            {
                headers: { 'Content-Type': 'application/json' },
                tags: { name: 'merge_pr' }
            }
        );

        const success = check(res, {
            'merge_pr: status 200 or 404': (r) => r.status === 200 || r.status === 404,
            'merge_pr: response time < 300ms': (r) => r.timings.duration < 300,
        });

        errorRate.add(!success);
        responseTime.add(res.timings.duration, { operation: 'merge_pr' });
    }

    // 20% — Reassign reviewer (POST /pullRequest/reassign)
    else {
        const res = http.post(
            `${BASE_URL}/pullRequest/reassign`,
            JSON.stringify({
                pull_request_id: prId,
                old_user_id: randomUser.user_id
            }),
            {
                headers: { 'Content-Type': 'application/json' },
                tags: { name: 'reassign' }
            }
        );

        const success = check(res, {
            'reassign: status 200, 404 or 409': (r) => r.status === 200 || r.status === 404 || r.status === 409,
            'reassign: response time < 300ms': (r) => r.timings.duration < 300,
        });

        errorRate.add(!success);
        responseTime.add(res.timings.duration, { operation: 'reassign' });
    }
}

export function teardown(data) {
    console.log('Test completed. Check the metrics above.');
}