import { Hex } from 'crypto-es/lib/core.js'
import { SHA1 } from 'crypto-es/lib/sha1'

const cron = async (event, env, ctx) => {

    // locale
    const locale = env.WORKERS_LOCALE

    try {
        const bingResponse = await (await fetch(`https://www.bing.com/HPImageArchive.aspx?idx=0&n=10&format=js&mkt=${locale}`)).json()
        if (Array.isArray(bingResponse?.images)) {
            const { results } = await env.DB.prepare("SELECT startdate FROM bing ORDER BY startdate DESC LIMIT 1;").all()
            //console.log(JSON.stringify(bingResponse.images, null, 4))
            const tmpList = bingResponse.images.filter(img => Number(results[0].startdate) < Number(img.startdate)).map(img => ({
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
                return 0
            }
            const stmt = env.DB.prepare("INSERT INTO bing (startdate, url, urlbase, copyright, copyrightlink, title, quiz, wp, hsh, drk, top, bot, hs) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
            await env.DB.batch(tmpList.map(img => stmt.bind(...Object.values(img))))
            console.log(`bing daily: ` + tmpList.map(img => [img.startdate, img.url].join(" -> ")).join(', '))

            // upload img
            const applicationKeyId = env.SECRET_WORKERS_APPLICATION_KEY_ID
            const applicationKey = env.SECRET_WORKERS_APPLICATION_KEY
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

            for (const img of tmpList) {

                const bingDailyImgBuffer = await (await fetch(`https://www.bing.com${img.urlbase}_UHD.jpg`)).arrayBuffer()

                //https://stackoverflow.com/questions/40031688/javascript-arraybuffer-to-hex
                const b2Upload = await (await fetch(b2UploadUrl.uploadUrl, {
                    headers: {
                        Authorization: b2UploadUrl.authorizationToken,
                        'Content-Type': 'image/jpeg',
                        'X-Bz-File-Name': encodeURIComponent(`bing/${img.startdate}.jpg`),
                        'Content-Length': bingDailyImgBuffer.byteLength,
                        'X-Bz-Content-Sha1': SHA1(Hex.parse([...new Uint8Array(bingDailyImgBuffer)]
                            .map(x => x.toString(16).padStart(2, '0'))
                            .join(''))).toString()
                    },
                    method: 'POST',
                    body: bingDailyImgBuffer
                })).json()
                //console.log(b2Upload)
                responseList.push(b2Upload)
            }
            console.log(JSON.stringify(responseList))


            return responseList
        }
    } catch (e) {
        // TODO push error message
        console.log(e)
    }
    return 0
}

export default cron