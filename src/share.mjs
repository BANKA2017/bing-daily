import { fileURLToPath } from 'node:url'
import { dirname } from 'node:path'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

const basePath = __dirname

const apiTemplate = (code = 403, message = 'Invalid Request', data = {}, version = 'online') => {
    if (version === 'v1') {
        return { error: code, message, data, version }
    } else {
        return { code, message, data, version }
    }
}

const TGPush = async (ALERT_TOKEN, ALERT_PUSH_TO, text = '') => {
    if (ALERT_TOKEN.length) {
        text = [...text]
        const partCount = Math.ceil(text.length / 3000)
        let tmpPartIndex = 0
        for (; tmpPartIndex < partCount; tmpPartIndex++) {
            try {
                const response = (await fetch(`https://api.telegram.org/bot${ALERT_TOKEN}/sendMessage`, {
                    headers: {
                        'content-type': 'application/json',
                    },
                    method: "POST",
                    body: JSON.stringify({
                        chat_id: ALERT_PUSH_TO,
                        text: text.slice(tmpPartIndex * 3000, tmpPartIndex * 3000 + 3000).join('')
                    })
                })).json()
                if (response?.ok) {
                    console.log(`TGPush: Successful to push log #part${tmpPartIndex} to chat ->${ALERT_PUSH_TO}<-`)
                } else {
                    console.log(`TGPush: Error #part${response?.description}`)
                }
            } catch (e) {
                console.log(e)
            }
        }
    }
}

export { apiTemplate, TGPush, basePath }
