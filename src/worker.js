import cron from "./cron.mjs"
import workersApi from "./router.mjs"
import { apiTemplate } from "./share.mjs"

export default {
    fetch: (req, env, ...args) =>
        workersApi
            .handle(req, env, ...args)
            .then((res) => {
                const WORKERS_ALLOW_ORIGIN = env.WORKERS_ALLOW_ORIGIN || []
                if (WORKERS_ALLOW_ORIGIN) {
                    const referer = req.headers.get('referer')
                    if (referer) {
                        const origin = new URL(referer).origin
                        const tmpReferer = WORKERS_ALLOW_ORIGIN.includes('*') ? '*' : WORKERS_ALLOW_ORIGIN.includes(origin) ? origin : ''
                        if (tmpReferer) {
                            res.headers.set('Access-Control-Allow-Origin', tmpReferer)
                        }
                    }
                }
                res.headers.set('X-Powered-By', 'Bing daily')
                res.headers.set('Access-Control-Allow-Methods', '*')
                res.headers.set('Access-Control-Allow-Credentials', 'true')
                return res
            })
            .catch((e) => {
                console.log(e)
                return new Response(JSON.stringify(apiTemplate(500, 'Unknown error', {}, 'global_api')), { status: 500 })
            }),
    async scheduled(event, env, ctx) {
        ctx.waitUntil(cron(event, env, ctx));
    },
}
