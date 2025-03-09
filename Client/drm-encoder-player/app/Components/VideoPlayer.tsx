export default function VideoPlayer() {
    return (
      <div className="p-4 border rounded-md">
        <h2 className="text-xl font-semibold mb-4">Test DRM Video Playback</h2>
        <video controls className="w-full">
          <source src="/path/to/encrypted/video.mpd" type="application/dash+xml" />
          Your browser does not support the video tag.
        </video>
      </div>
    );
  }