/* global videojs */
(function () {
  function getQuery(k) {
    const p = new URLSearchParams(location.search);
    return p.get(k) || "";
  }
  // The backend serves this HTML at /files/<path> when media extension and no raw=1
  // We need to reconstruct the actual media file URL so we can stream it directly with range requests.
  // Original requested path is after /files/ in the pathname.
  function extractMediaPath() {
    const idx = location.pathname.indexOf("/files/");
    if (idx === -1) return "";
    // Keep everything after /files/
    return decodeURIComponent(location.pathname.substring(idx + 7));
  }

  const relPath = extractMediaPath();
  const fileTitleEl = document.getElementById("fileTitle");
  fileTitleEl.textContent = relPath || "Media";

  const mediaUrl =
    "/files/" + relPath + (relPath.includes("?") ? "&" : "?") + "raw=1";

  // Decide if audio-only
  const lower = relPath.toLowerCase();
  const isAudio = /\.(mp3|wav|ogg|m4a|flac|aac)$/.test(lower);

  const videoEl = document.getElementById("videoPlayer");
  if (isAudio) {
    videoEl.classList.add("vjs-audio");
    videoEl.setAttribute("style", "max-width:720px; width:100%;");
  } else {
    videoEl.setAttribute("playsinline", "");
  }

  // Provide source – rely on server Content-Type
  const sourceEl = document.createElement("source");
  sourceEl.src = mediaUrl;
  sourceEl.type = guessType(lower);
  videoEl.appendChild(sourceEl);

  const player = videojs(videoEl, {
    controls: true,
    preload: "auto",
    playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 2],
    fluid: !isAudio,
    userActions: { hotkeys: false },
    controlBar: {
      volumePanel: { inline: false, vertical: true },
    },
  });

  // Keyboard shortcuts: Left / Right = 5s seek, Space = toggle, F = fullscreen
  document.addEventListener("keydown", (e) => {
    // Ignore if focused inside input/textarea/contentEditable
    const tag = document.activeElement && document.activeElement.tagName;
    if (
      tag === "INPUT" ||
      tag === "TEXTAREA" ||
      (document.activeElement && document.activeElement.isContentEditable)
    )
      return;

    const seekStep = 5; // seconds
    switch (e.key) {
      case "ArrowLeft":
        e.preventDefault();
        player.currentTime(Math.max(0, player.currentTime() - seekStep));
        flashSeek("-" + seekStep + "s");
        break;
      case "ArrowRight":
        e.preventDefault();
        player.currentTime(
          Math.min(
            player.duration() || Infinity,
            player.currentTime() + seekStep
          )
        );
        flashSeek("+" + seekStep + "s");
        break;
      case " ": // Space
        e.preventDefault();
        if (player.paused()) player.play();
        else player.pause();
        break;
      case "f":
      case "F":
        e.preventDefault();
        if (player.isFullscreen()) player.exitFullscreen();
        else player.requestFullscreen();
        break;
      default:
        return;
    }
  });

  function flashSeek(text) {
    let el = document.getElementById("seekFlash");
    if (!el) {
      el = document.createElement("div");
      el.id = "seekFlash";
      Object.assign(el.style, {
        position: "fixed",
        left: "50%",
        top: "20%",
        transform: "translateX(-50%)",
        background: "rgba(0,0,0,.6)",
        color: "#fff",
        padding: "6px 12px",
        borderRadius: "20px",
        fontSize: "14px",
        fontWeight: "500",
        zIndex: 9999,
        opacity: "0",
        transition: "opacity .15s",
      });
      document.body.appendChild(el);
    }
    el.textContent = text;
    el.style.opacity = "1";
    clearTimeout(flashSeek._t);
    flashSeek._t = setTimeout(() => {
      el.style.opacity = "0";
    }, 350);
  }

  function guessType(p) {
    if (p.endsWith(".mp4")) return "video/mp4";
    if (p.endsWith(".webm")) return "video/webm";
    if (p.endsWith(".ogv") || p.endsWith(".ogg")) return "video/ogg";
    if (p.endsWith(".mp3")) return "audio/mpeg";
    if (p.endsWith(".wav")) return "audio/wav";
    if (p.endsWith(".m4a")) return "audio/mp4";
    if (p.endsWith(".flac")) return "audio/flac";
    if (p.endsWith(".aac")) return "audio/aac";
    return "";
  }
})();
