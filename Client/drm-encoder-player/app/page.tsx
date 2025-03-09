'use client'
import { useState } from 'react';
import Head from 'next/head';
import EncryptVideo from './Components/EncryptVideo';
import TestPlayback from './Components/TestPLayback';

const Home = () => {
  const [selectedTab, setSelectedTab] = useState('encrypt');
  const [encryptionType, setEncryptionType] = useState('Widevine');
  const [useTestServer, setUseTestServer] = useState(false);

  return (
    <div className="min-h-screen bg-gray-100 text-gray-900">
      <Head>
        <title>Video Encryption</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <header className="bg-white shadow">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-xl font-bold">Video Encryption</h1>
          <nav>
            <button
              className={`px-4 py-2 mx-2 ${selectedTab === 'encrypt' ? 'bg-blue-500 text-white' : 'bg-gray-200'}`}
              onClick={() => setSelectedTab('encrypt')}
            >
              Encrypt Video
            </button>
            <button
              className={`px-4 py-2 mx-2 ${selectedTab === 'playback' ? 'bg-blue-500 text-white' : 'bg-gray-200'}`}
              onClick={() => setSelectedTab('playback')}
            >
              Test Playback
            </button>
          </nav>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        {selectedTab === 'encrypt' ? (
          <EncryptVideo
            encryptionType={encryptionType}
            setEncryptionType={setEncryptionType}
            useTestServer={useTestServer}
            setUseTestServer={setUseTestServer}
          />
        ) : (
          <TestPlayback
            encryptionType={encryptionType}
            setEncryptionType={setEncryptionType}
            useTestServer={useTestServer}
            setUseTestServer={setUseTestServer}
          />
        )}
      </main>
    </div>
  );
};
export default Home