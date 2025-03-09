import EncodingConfig from "./EncodingConfig";

const TestPlayback = ({ encryptionType, setEncryptionType, useTestServer, setUseTestServer }: any) => {
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
                <div className="bg-white p-6 rounded shadow">
                    <div className="mt-6">
                        <video className="w-full" controls>
                            <source src="path/to/your/video.mp4" type="video/mp4" />
                            Your browser does not support the video tag.
                        </video>
                    </div>
                </div>
            </section>
        </>
    );
};
export default TestPlayback