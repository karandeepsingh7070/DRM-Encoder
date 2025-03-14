'use client'
import React, { useEffect } from 'react'

type Config = {
    encryptionType: string
    setEncryptionType: (type: string) => void
    useTestServer: boolean
    setUseTestServer: (type: boolean) => void,
    includeAudio: string,
    segmentSize: string,
    setSegmentSize: (val:string) => void,
    setIncludeAudio: (val:string) => void,
}
const EncodingConfig = ({ encryptionType, setEncryptionType, useTestServer, setUseTestServer, segmentSize, includeAudio, setSegmentSize,setIncludeAudio }: Config) => {

    const saveToSession = (e: React.ChangeEvent<HTMLInputElement>, type: string) => {
        sessionStorage.setItem(type, e.target.value)
    }
    useEffect(() => {
        setUseTestServer(true)
    }, [])
    return (<>
        <section className="bg-white p-6 rounded shadow">
            <h3 className="text-xl font-semibold mb-4">Encryption Settings</h3>
            <div className="flex">
                <label className="block mb-4">
                    <span className="text-gray-700">Encryption Type</span>
                    <select
                        className="block w-full mt-1 bg-blue-500 text-white p-1.5"
                        value={encryptionType}
                        onChange={(e) => setEncryptionType(e.target.value)}
                    >
                        <option value="Widevine">Widevine</option>
                        <option disabled={true} value="Playready">PlayReady</option>
                        <option value="RawKey">Raw Key</option>
                    </select>
                </label>
            </div>
            {encryptionType === 'Widevine' && (
                <>
                    <div className="flex gap-2">
                        <label className="block mb-4">
                            <span className="text-gray-700">License Server URL</span>
                            <input
                                type="text"
                                onChange={(e) => saveToSession(e, "widevine-server")}
                                className={`${useTestServer ? "bg-gray-400" : "bg-blue-300"} block w-full mt-1  text-black p-1.5`}
                                disabled={useTestServer}
                            />
                        </label>
                        <label className="block mb-4">
                            <span className="text-gray-700">AES Signing Key</span>
                            <input
                                type="text"
                                onChange={(e) => saveToSession(e, "AES-Sign-key")}
                                className={`${useTestServer ? "bg-gray-400" : "bg-blue-300"} block w-full mt-1  text-black p-1.5`}
                                disabled={useTestServer}
                            />
                        </label>
                    </div>
                    <div className="flex gap-2">
                        <label className="block mb-4">
                            <span className="text-gray-700">AES Signing IV</span>
                            <input
                                type="text"
                                onChange={(e) => saveToSession(e, "AES-Sign-4")}
                                className={`${useTestServer ? "bg-gray-400" : "bg-blue-300"} block w-full mt-1  text-black p-1.5`}
                                disabled={useTestServer}
                            />
                        </label>
                        <label className="flex items-center mb-4">
                            <input
                                type="checkbox"
                                className="mr-2 bg-blue-300 text-white p-1.5"
                                checked={useTestServer}
                                onChange={() => setUseTestServer(!useTestServer)}
                            />
                            <span>Use Test Widevine Server</span>
                        </label>
                    </div>
                </>
            )}
            {encryptionType === 'RawKey' && (
                <div className="flex gap-2">
                    <label className="block mb-4">
                        <span className="text-gray-700">Key ID</span>
                        <input type="text" className="block w-full mt-1 bg-blue-300 text-white p-1.5"
                            value={"07507c220e89a23e20b25a2d03b74d53"}
                            onChange={(e) => saveToSession(e, "keyId")}
                        />
                    </label>
                    <label className="block mb-4">
                        <span className="text-gray-700">Value</span>
                        <input type="text" className="block w-full mt-1 bg-blue-300 text-white p-1.5"
                            value={"6e19d3fabf454e4f0be778844354cf81"}
                            onChange={(e) => saveToSession(e, "val")} />
                    </label>
                </div>
            )}
            {includeAudio && <div className="flex gap-2">
                <label className="flex gap-1 mb-4">
                    <span className="text-gray-700">Enable Audio</span>
                    <select name="enableAudio" id="enableAudio" className='bg-blue-500 text-white' value={includeAudio} 
                    onChange={(e) => setIncludeAudio(e.target.value)}>
                        <option className='p-1' value="yes">Enabled</option>
                        <option value="no">Disabled</option>
                    </select>
                </label>
            </div>}
            {segmentSize && <div className="flex flex-nowrap gap-2 ">
            <label className="flex gap-1 mb-4 align-middle">
                    <span className="text-gray-700 w-40 flex">Add segment size</span>
                    <input type="number" pattern="^[1-9]$" maxLength={1} className="block w-30 mt-1 bg-blue-300 text-white p-1.5"
                            value={segmentSize}
                            onChange={(e) => setSegmentSize(e.target.value)} />
                </label>
            </div>}
        </section>
    </>)
}

export default EncodingConfig