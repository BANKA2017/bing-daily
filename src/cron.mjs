import { Hex } from 'crypto-es/lib/core.js'
import { SHA1 } from 'crypto-es/lib/sha1.js'
//import { TGPush } from './share.mjs'
import { FastAverageColor } from 'fast-average-color'
import inkjet from 'inkjet'
import { encode } from 'blurhash'
import Database from 'better-sqlite3';
import { basePath } from './share.mjs'
import 'dotenv/config'

const db = new Database(basePath + '/../db/bing_kv.sqlite3');
//db.pragma('journal_mode = WAL');

const env = process.env

// locale
const locale = env.WORKERS_LOCALE

try {
    const bingResponse = await (await fetch(`https://www.bing.com/HPImageArchive.aspx?idx=0&n=10&format=js&mkt=${locale}`)).json()
    if (Array.isArray(bingResponse?.images)) {
        const results = db.prepare("SELECT startdate FROM bing ORDER BY startdate DESC LIMIT 1;").all()
        //console.log(JSON.stringify(bingResponse.images, null, 4))
        let tmpList = bingResponse.images.filter(img => Number(results[0].startdate) < Number(img.startdate)).map(img => ({
            startdate: img.startdate,
            url: img.url,
            urlbase: img.urlbase,
            copyright: img.copyright,
            copyrightlink: img.copyrightlink,
            title: img.title,
            quiz: img.quiz,
            wp: img.wp ? 1 : 0,
            hsh: img.hsh,
            drk: img.drk,
            top: img.top,
            bot: img.bot,
            hs: JSON.stringify(img.hs)
        }))
        if ((tmpList || []).length <= 0) {
            process.exit()
        }
        // upload img
        const applicationKeyId = env.B2_APPLICATION_KEY_ID
        const applicationKey = env.B2_APPLICATION_KEY
        //get account
        const b2AuthorizeAccount = await (await fetch('https://api.backblazeb2.com/b2api/v2/b2_authorize_account', {
            headers: {
                Authorization: 'Basic ' + btoa(`${applicationKeyId}:${applicationKey}`)
            }
        })).json()
        const b2UploadUrl = await (await fetch(`${b2AuthorizeAccount.apiUrl}/b2api/v2/b2_get_upload_url?bucketId=${b2AuthorizeAccount.allowed.bucketId}`, {
            headers: {
                Authorization: b2AuthorizeAccount.authorizationToken
            }
        })).json()
        const responseList = []

        for (const index in tmpList) {
            const img = tmpList[index]

            const bingDailyImgBuffer = await (await fetch(`https://www.bing.com${img.urlbase}_UHD.jpg`)).arrayBuffer()
            const bingDailyImgSmallBuffer = await (await fetch(`https://www.bing.com${img.urlbase}_UHD.jpg&rf=LaDigue_UHD.jpg&pid=hp&w=256&h=128&rs=1&c=4`)).arrayBuffer()

            //https://stackoverflow.com/questions/40031688/javascript-arraybuffer-to-hex
            const b2Upload = await (await fetch(b2UploadUrl.uploadUrl, {
                headers: {
                    Authorization: b2UploadUrl.authorizationToken,
                    'Content-Type': 'image/jpeg',
                    'X-Bz-File-Name': encodeURIComponent(`bing/${img.startdate}.jpg`),
                    'Content-Length': bingDailyImgBuffer.byteLength,
                    'X-Bz-Content-Sha1': SHA1(Hex.parse((bingDailyImgBuffer => {
                        let buf = new Uint8Array(bingDailyImgBuffer)
                        let tmpArray = new Array(buf.byteLength)
                        for (let x = 0; x < buf.byteLength; x++) {
                            tmpArray[x] = buf[x].toString(16).padStart(2, "0")
                        }
                        return tmpArray.join('')
                    })(bingDailyImgBuffer))).toString()
                },
                method: 'POST',
                body: bingDailyImgBuffer
            })).json()
            //console.log(b2Upload)
            responseList.push(b2Upload)
            const fac = new FastAverageColor();
            let color, blurhash, width, height

            //small
            inkjet.decode(bingDailyImgSmallBuffer, (err, decoded) => {
                //console.log(err, decoded)
                // decoded: { width: number, height: number, data: Uint8Array }
                color = fac.getColorFromArray4(new Uint8ClampedArray(decoded.data), { step: 5 });
                blurhash = encode(new Uint8ClampedArray(decoded.data), decoded.width, decoded.height, 4, 4)
            })
            //uhd
            inkjet.decode(bingDailyImgBuffer, (err, decoded) => {
                //console.log(err, decoded)
                width = decoded.width
                height = decoded.height
            })
            tmpList[index].color = color[0].toString(16).padStart(2, '0') + color[1].toString(16).padStart(2, '0') + color[2].toString(16).padStart(2, '0')
            tmpList[index].blurhash = blurhash
            tmpList[index].width = width
            tmpList[index].height = height

            //TGPush()
        }
        //console.log(JSON.stringify(tmpList))
        console.log(JSON.stringify(responseList))
        const stmt = db.prepare("INSERT INTO bing (startdate, url, urlbase, copyright, copyrightlink, title, quiz, wp, hsh, drk, top, bot, hs, color, blurhash, width, height) values (@startdate, @url, @urlbase, @copyright, @copyrightlink, @title, @quiz, @wp, @hsh, @drk, @top, @bot, @hs, @color, @blurhash, @width, @height)")
        db.transaction((imgs) => {
            for (const img of imgs) {stmt.run(img)}
        })(tmpList)
        //db.batch(tmpList.map(img => stmt.bind(...Object.values(img))))
        console.log(`bing daily: ` + tmpList.map(img => [img.startdate, img.url].join(" -> ")).join(', '))
    }
} catch (e) {
    // TODO push error message
    console.log(e)
}
