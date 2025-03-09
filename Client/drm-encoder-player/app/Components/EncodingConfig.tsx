import React from 'react'

const EncodingConfig = ({ encryptionType, setEncryptionType, useTestServer, setUseTestServer }: any) => {
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
                        <option value="Playready">PlayReady</option>
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
                                className={`${useTestServer ? "bg-gray-400" : "bg-blue-300"} block w-full mt-1  text-black p-1.5`}
                                disabled={useTestServer}
                            />
                        </label>
                        <label className="block mb-4">
                            <span className="text-gray-700">AES Signing Key</span>
                            <input
                                type="text"
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
                        <input type="text" className="block w-full mt-1 bg-blue-300 text-white p-1.5" />
                    </label>
                    <label className="block mb-4">
                        <span className="text-gray-700">Value</span>
                        <input type="text" className="block w-full mt-1 bg-blue-300 text-white p-1.5" />
                    </label>
                </div>
            )}
        </section>
    </>)
}

export default EncodingConfig