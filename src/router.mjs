import { apiTemplate } from './share.mjs'
import express from 'express'
import { basePath } from './share.mjs'
import Database from 'better-sqlite3';
const app = express()

const db = new Database(basePath + '/../db/bing_kv.sqlite3');
const EXPRESS_ALLOW_ORIGIN = ["*"]

app.use((req, res, next) => {
    
    if (EXPRESS_ALLOW_ORIGIN && req.headers.referer) {
        const origin = new URL(req.headers.referer).origin
        const tmpReferer = EXPRESS_ALLOW_ORIGIN.includes('*') ? '*' : EXPRESS_ALLOW_ORIGIN.includes(origin) ? origin : ''
        if (tmpReferer) {
            res.append('Access-Control-Allow-Origin', tmpReferer)
        }
    }
    res.setHeader('X-Powered-By', 'Bing daily')
    res.setHeader('Access-Control-Allow-Methods', '*')
    res.setHeader('Access-Control-Allow-Credentials', 'true')
    next()
})

//favicon
app.all('/favicon.ico', (req, res) => {res.json(null)})

//robots.txt
app.all('/robots.txt', (req, res) => {
    res.send('User-agent: *\nDisallow: /*')
})

// DO NOT UNCOMMENT, THE RESPONSE WILL LEAK B2 FILE INFO
//workersApi.get('/test/upload/run', async (req, env) => {
//    const uploadData = await cron(null, env, null)
//    return env.json(apiTemplate(200, 'OK', [], 'online'), 200)
//})

app.get('/v1/data/list/', (req, res) => {
    //count
    let count = Number(req.query.count) || 16
    if (count < 1) { count = 1 }
    if (count > 100) { count = 100 }

    //date
    let date = Number(req.query.date) || 30000101 // AD 3000-01-01

    const results = db.prepare("SELECT startdate, url, urlbase, copyright, copyrightlink, title, quiz, blurhash, color, width, height FROM bing WHERE startdate < ? ORDER BY startdate DESC LIMIT ?;").all(date, count + 1)
    const more = results.length === count + 1
    //console.log(results, date)
    res.json(apiTemplate(200, 'OK', {images: results.slice(0, more ? -1 : undefined), more}, 'online'))
})

app.all('*', (req, res) => res.status(403).json(apiTemplate(403, 'Invalid Request', {}, 'global_api')))

app.use((err, req, res, next) => {
    console.log(new Date(), err)
    res.status(500).json(apiTemplate(500, 'Unknown error', {}, 'global_api'))
})

app.listen(3000, () => {
    console.log(`V3Api listening on port 3000`)
})
