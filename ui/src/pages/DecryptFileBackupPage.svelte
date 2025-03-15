<script lang="ts">
  import { onMount } from "svelte";
  import { zxing } from "../zxing";

  let passphrase = $state("");
  let message = $state("");
  let result = $state("");

  onMount(() => {
    const resultElement = document.getElementById("result");
    const canvas = document.getElementById("qr-reader") as HTMLCanvasElement;
    const ctx = canvas.getContext("2d", { willReadFrequently: true });
    const video = document.createElement("video");
    video.setAttribute("id", "video");
    video.setAttribute("width", canvas.width);
    video.setAttribute("height", canvas.height);
    video.setAttribute("autoplay", "");

    function readBarcodeFromCanvas(canvas) {
      var imgWidth = canvas.width;
      var imgHeight = canvas.height;
      var imageData = canvas
        .getContext("2d")
        .getImageData(0, 0, imgWidth, imgHeight);
      var sourceBuffer = imageData.data;

      if (zxing != null) {
        var buffer = zxing._malloc(sourceBuffer.byteLength);
        zxing.HEAPU8.set(sourceBuffer, buffer);
        var result = zxing.readBarcodeFromPixmap(
          buffer,
          imgWidth,
          imgHeight,
          true,
          "QRCode",
        );
        zxing._free(buffer);
        return result;
      } else {
        return { error: "ZXing not yet initialized" };
      }
    }

    function drawResult(code) {
      ctx.beginPath();
      ctx.lineWidth = 4;
      ctx.strokeStyle = "red";
      // ctx.textAlign = "center";
      // ctx.fillStyle = "#green"
      // ctx.font = "25px Arial";
      // ctx.fontWeight = "bold";
      ctx.moveTo(code.position.topLeft.x, code.position.topLeft.y);
      ctx.lineTo(code.position.topRight.x, code.position.topRight.y);
      ctx.lineTo(code.position.bottomRight.x, code.position.bottomRight.y);
      ctx.lineTo(code.position.bottomLeft.x, code.position.bottomLeft.y);
      ctx.lineTo(code.position.topLeft.x, code.position.topLeft.y);
      ctx.stroke();
      // ctx.fillText(code.text, (topLeft.x + bottomRight.x) / 2, (topLeft.y + bottomRight.y) / 2);
    }

    function escapeTags(htmlStr) {
      return htmlStr
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#39;");
    }

    const processFrame = function () {
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height);

      const code = readBarcodeFromCanvas(canvas);
      if (code.format) {
        console.log(code);
        resultElement.innerText = code.format + ": " + code.bytes;
        drawResult(code);
      } else {
        resultElement.innerText = "No barcode found";
      }
      requestAnimationFrame(processFrame);
    };

    const updateVideoStream = function () {
      // To ensure the camera switch, it is advisable to free up the media resources
      if (video.srcObject)
        video.srcObject.getTracks().forEach((track) => track.stop());

      navigator.mediaDevices
        .getUserMedia({ video: { facingMode: "" }, audio: false })
        .then(function (stream) {
          video.srcObject = stream;
          video.setAttribute("playsinline", true); // required to tell iOS safari we don't want fullscreen
          video.play();
          processFrame();
        })
        .catch(function (error) {
          console.error("Error accessing camera:", error);
        });
    };

    updateVideoStream();
  });
</script>

<div>
  Decrypt a file backup.

  <div>
    <label
      >Enter a passphrase
      <input bind:value={passphrase} name="passphrase" type="text" /></label
    >
  </div>

  <canvas id="qr-reader" width="640" height="480"></canvas>

  <div id="result"></div>

  <div>{message}</div>

  <div>{result}</div>
</div>
