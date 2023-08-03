import { Router } from 'itty-router'
import { PostBodyParser, json, xml, apiTemplate } from './share.mjs'
//import cron from './cron.mjs'

const workersApi = Router()

workersApi.all('*', (req, env) => {
    env.json = json
    env.xml = xml
    env.PostBodyParser = PostBodyParser
    req.cookies = Object.fromEntries(
        (req.headers.get('cookie') || '')
            .split(';')
            .map((cookie) => cookie.trim().split('='))
            .filter((cookie) => cookie.length === 2)
    )
})

//favicon
workersApi.all('/favicon.ico', () => new Response(null, { status: 200 }))

//robots.txt
workersApi.all('/robots.txt', () => new Response('User-agent: *\nDisallow: /*', { status: 200 }))

// DO NOT UNCOMMENT, THE RESPONSE WILL LEAK B2 FILE INFO
//workersApi.get('/test/upload/run', async (req, env) => {
//    const uploadData = await cron(null, env, null)
//    return env.json(apiTemplate(200, 'OK', [], 'online'), 200)
//})

workersApi.get('/v1/data/list/', async (req, env) => {
    //count
    let count = Number(req.query.count) || 1
    if (count < 1) { count = 1 }
    if (count > 10) { count = 10 }

    //date
    let date = Number(req.query.date) || 30000101 // AD 3000-01-01

    const { results } = await env.DB.prepare("SELECT startdate, url, urlbase, copyright, copyrightlink, title, quiz, blurhash, color, width, height FROM bing WHERE startdate < ?2 ORDER BY startdate DESC LIMIT ?1;").bind(count, date).all()
    //console.log(results, date)
    return env.json(apiTemplate(200, 'OK', results, 'online'), 200)
})

workersApi.all('*', () => new Response(JSON.stringify(apiTemplate(403, 'Invalid Request', {}, 'global_api')), { status: 403 }))


export default workersApi