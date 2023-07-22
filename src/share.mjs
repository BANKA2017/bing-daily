const json = (data, status = 200) =>
    new Response(JSON.stringify(data), {
        status,
        headers: {
            'content-type': 'application/json'
        }
    })

const xml = (data, status = 200) =>
    new Response(data, {
        status,
        headers: {
            'content-type': 'application/xml;charset=UTF-8'
        }
    })

const PostBodyParser = async (req, defaultValue = new Map([])) => {
    if (req.body) {
        const reader = req.body.getReader()
        const pipe = []
        while (true) {
            const { done, value } = await reader.read()
            if (done) {
                break
            }
            pipe.push(value)
        }
        //https://gist.github.com/72lions/4528834
        let offset = 0
        let body = new Uint8Array(pipe.reduce((acc, cur) => acc + cur.byteLength, 0))
        for (const chunk of pipe) {
            body.set(new Uint8Array(chunk), offset)
            offset += chunk.byteLength
        }
        //TODO json parser
        req.postBody = new URLSearchParams(new TextDecoder('utf-8').decode(body))
    } else {
        return defaultValue
    }
}

const apiTemplate = (code = 403, message = 'Invalid Request', data = {}, version = 'online') => {
    if (version === 'v1') {
        return { error: code, message, data, version }
    } else {
        return { code, message, data, version }
    }
}

export { json, xml, PostBodyParser, apiTemplate }
