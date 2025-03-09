'use client'
import { useRef } from "react";
import EncodingConfig from "./EncodingConfig";
import shaka from "shaka-player";

const TestPlayback = ({ encryptionType, setEncryptionType, useTestServer, setUseTestServer, manifestUrl }: any) => {
    const videoRef = useRef<HTMLVideoElement>(null);

    async function initPlayer() {
        if (!videoRef.current) return;

        shaka.polyfill.installAll(); // Install polyfills for better compatibility

        const player = new shaka.Player(videoRef.current);

        player.addEventListener("error", (e) => {
            console.error("Error loading Shaka Player:", e);
        });

        try {
            await player.configure({
                drm: {
                    clearKeys: {
                        "07507c220e89a23e20b25a2d03b74d53": "6e19d3fabf454e4f0be778844354cf81"
                    }
                }
            });
            await player.load("http://localhost:8080/uploads/stream.mpd");
            console.log("The video has loaded successfully!");
        } catch (error) {
            console.error("Error loading manifest:", error);
        }
    }
    return (
        <>
            <section className="mb-8">
                <h2 className="text-2xl font-bold mb-4">Test Playback</h2>
                <EncodingConfig
                    encryptionType={encryptionType}
                    setEncryptionType={setEncryptionType}
                    useTestServer={useTestServer}
                    setUseTestServer={setUseTestServer}
                />
            </section>
            <section className="mb-8 bg-white p-6 rounded shadow">
                <button className="px-4 py-2 cursor-pointer bg-blue-500 text-white rounded" onClick={initPlayer}>Play Video</button>
                <h6 className="text-red-500 pt-2 pb-3">IMPORTANT : Keep the encryption type & keys same used during encoding</h6>
                <div className="w-full flex flex-col items-center">
                    <video
                        ref={videoRef}
                        className="w-full max-w-3xl border-2 border-gray-500 rounded-lg aspect-16/9"
                        controls
                    />
                </div>
            </section>
        </>
    );
};
export default TestPlayback