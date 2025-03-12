'use client'
import { useState } from "react";
import EncodingConfig from "./EncodingConfig";

type Config = {
    encryptionType: string
    setEncryptionType: (type:string) => void
    useTestServer: boolean
    setUseTestServer: (type:boolean) => void
    segmentSize: string,
    includeAudio: string,
    setSegmentSize: (val:string) => void,
    setIncludeAudio: (val:string) => void,
}
const EncryptVideo = ({ encryptionType, setEncryptionType, useTestServer, setUseTestServer, segmentSize, includeAudio,setSegmentSize,setIncludeAudio }: Config) => {

    console.log(segmentSize, includeAudio)

    const [file, setFile] = useState<File | null>(null);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState("");

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            setFile(e.target.files[0]);
        }
    };
    const handleUpload = async () => {
        if (!file) {
            setMessage("Please select a file first.");
            return;
        }

        setLoading(true);
        setMessage("");

        const formData = new FormData();
        formData.append("video", file);
        formData.append("encryptionType", encryptionType);
        // formData.append("segmentSize", segmentSize);
        formData.append("includeAudio", includeAudio);

        try {
            const response = await fetch("http://localhost:8080/upload", {
                method: "POST",
                body: formData,
            });

            const data = await response.text();
            if (response.ok) {
                setMessage(data);
            } else {
                setMessage("Upload failed: " + data);
            }
        } catch (error) {
            setMessage("Error uploading file.");
            console.error("Upload error:", error);
        } finally {
            setLoading(false);
        }
    };

    const fetchEncryptedFiles = () => {
        const link = document.createElement("a");
        link.href = "http://localhost:8080/get-files";
        link.setAttribute("download", "encrypted_videos.zip");
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }
    return (
        <>
            <section className="mb-8">
                <h2 className="text-2xl font-bold mb-4">Encrypt Video</h2>
                <EncodingConfig
                    encryptionType={encryptionType}
                    setEncryptionType={setEncryptionType}
                    useTestServer={useTestServer}
                    setUseTestServer={setUseTestServer}
                    segmentSize={segmentSize}
                    includeAudio={includeAudio}
                    setSegmentSize={setSegmentSize}
                    setIncludeAudio={setIncludeAudio}
                />
                <div className="bg-white p-6 rounded shadow">
                    <label className="block mb-4 p-5 border-2">
                        <span className="text-gray-700">{loading ? "Uploading..." : "Upload Video"}</span>
                        <input type="file" accept="video/mp4" className=" w-full mt-1 flex justify-center cursor-pointer"
                            onChange={handleFileChange}
                        />
                    </label>
                    <button disabled={!file || loading} className="px-4 py-2 cursor-pointer bg-blue-500 text-white rounded disabled:opacity-50" onClick={handleUpload}>Encrypt Video</button>
                    {message && <p className="mt-2">{message}</p>}
                    <div className="mt-6">
                        <ul className="list-disc list-inside">
                            <li>Convert Video to DRM-Protected Format</li>
                            <li>Secure your video content with Widevine, PlayReady, or Raw Key encryption.</li>
                            <li>Upload & Configure Your Video</li>
                            <li>Customize your video encoding settings before encryption.</li>
                            <li>Test Your Encrypted Video</li>
                            <li>Ensure playback compatibility with different DRM types.</li>
                        </ul>
                    </div>
                    <button disabled={!file || loading} className="px-4 py-2 cursor-pointer mt-7 bg-blue-500 text-white rounded disabled:opacity-50" onClick={fetchEncryptedFiles}>Download Encrypted Segments</button>
                </div>
            </section>
        </>
    );
};
export default EncryptVideo