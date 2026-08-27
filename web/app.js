"use strict";

// app.js 只负责收集输入、调用后端 API 并渲染返回的数值，
// 不含任何扩散求解逻辑。

const $ = (id) => document.getElementById(id);
const fields = {
  q: () => parseFloat($("f-q").value),
  u: () => parseFloat($("f-u").value),
  stab: () => $("f-stab").value,
  x: () => parseFloat($("f-x").value),
  y: () => parseFloat($("f-y").value),
  z: () => parseFloat($("f-z").value),
  x0: () => parseFloat($("f-x0").value),
  x1: () => parseFloat($("f-x1").value),
  n: () => parseInt($("f-n").value, 10),
  mode: () => {
    const active = document.querySelector("#sourceMode button.active");
    return active ? active.dataset.mode : "height";
  },
};

function buildSource() {
  if (fields.mode() === "stack") {
    return {
      stack: {
        height: parseFloat($("f-hs").value),
        exit_velocity: parseFloat($("f-ws").value),
        radius: parseFloat($("f-rs").value),
        gas_temperature: parseFloat($("f-ts").value),
        ambient_temperature: parseFloat($("f-ta").value),
        inversion_top: parseFloat($("f-zinv").value),
      },
    };
  }
  return { height: parseFloat($("f-height").value) };
}

function setError(msg) {
  const el = $("error");
  if (msg) {
    el.textContent = msg;
    el.style.display = "block";
  } else {
    el.style.display = "none";
  }
}

function setBusy(busy) {
  $("btn-conc").disabled = busy;
  $("btn-axis").disabled = busy;
}

async function api(path, payload) {
  const res = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(data.error || ("HTTP " + res.status));
  }
  return data;
}

function renderPoint(p) {
  const metric = (label, value) =>
    `<div class="cell"><b>${value}</b><span>${label}</span></div>`;
  $("result").innerHTML = `
    <div class="metric">
      ${metric("地面/受体浓度 (g/m³)", p.concentration.toExponential(4))}
      ${metric("有效源高 H (m)", p.effective_height.toFixed(2))}
      ${metric("σy (m)", p.sigma_y.toFixed(2))}
      ${metric("σz (m)", p.sigma_z.toFixed(2))}
      ${metric("抬升 ΔH (m)", p.plume_rise.toFixed(2))}
      ${metric("逆温截断", p.capped ? "是" : "否")}
    </div>`;
}

function renderAxis(a) {
  const pts = a.points;
  const metric = (label, value) =>
    `<div class="cell"><b>${value}</b><span>${label}</span></div>`;
  const rows = pts
    .slice(0, 200)
    .map(
      (p) =>
        `<tr><td>${p.x.toFixed(1)}</td><td>${p.concentration.toExponential(4)}</td></tr>`
    )
    .join("");
  $("result").innerHTML = `
    <div class="metric">
      ${metric("点数", pts.length)}
      ${metric("有效源高 H (m)", a.effective_height.toFixed(2))}
      ${metric("抬升 ΔH (m)", a.plume_rise.toFixed(2))}
      ${metric("逆温截断", a.capped ? "是" : "否")}
    </div>
    <h2>下风向曲线（来自 /api/axis）</h2>
    <table>
      <tr><th>x (m)</th><th>浓度 (g/m³)</th></tr>
      ${rows}
    </table>`;
  drawChart(pts);
}

function drawChart(pts) {
  const W = 720;
  const H = 280;
  const padL = 70;
  const padR = 20;
  const padT = 16;
  const padB = 40;
  const xMin = pts[0].x;
  const xMax = pts[pts.length - 1].x;
  let cMax = 0;
  for (const p of pts) {
    if (p.concentration > cMax) cMax = p.concentration;
  }
  const yMin = 1e-12;
  const xSpan = xMax - xMin || 1;
  const sx = (x) => padL + ((x - xMin) / xSpan) * (W - padL - padR);
  const sy = (c) => {
    const t = Math.log(Math.max(c, yMin) / yMin) / Math.log(cMax / yMin || 1);
    return padT + (1 - t) * (H - padT - padB);
  };

  let path = "";
  for (const p of pts) {
    path += (path ? "L" : "M") + sx(p.x).toFixed(1) + "," + sy(p.concentration).toFixed(1);
  }
  let ticks = "";
  const nTicks = 6;
  for (let i = 0; i < nTicks; i++) {
    const f = yMin * Math.pow(cMax / yMin, i / (nTicks - 1));
    const y = sy(f);
    ticks +=
      `<text x="${padL - 8}" y="${(y + 4).toFixed(1)}" text-anchor="end" font-size="11" fill="#6b7684">${f.toExponential(1)}</text>` +
      `<line x1="${padL}" y1="${y.toFixed(1)}" x2="${W - padR}" y2="${y.toFixed(1)}" stroke="#dde2e8" stroke-width="1"/>`;
  }
  const xTicks = 5;
  let xLabels = "";
  for (let i = 0; i < xTicks; i++) {
    const x = xMin + ((xMax - xMin) * i) / (xTicks - 1);
    xLabels +=
      `<text x="${sx(x).toFixed(1)}" y="${H - 20}" text-anchor="middle" font-size="11" fill="#6b7684">${x.toFixed(0)}</text>`;
  }

  $("chart").innerHTML = `
    <svg width="${W}" height="${H}" viewBox="0 0 ${W} ${H}" role="img" aria-label="下风向浓度曲线">
      <line x1="${padL}" y1="${padT}" x2="${padL}" y2="${H - padB}" stroke="#1f2733"/>
      <line x1="${padL}" y1="${H - padB}" x2="${W - padR}" y2="${H - padB}" stroke="#1f2733"/>
      ${ticks}
      ${xLabels}
      <path d="${path}" fill="none" stroke="#2f6fed" stroke-width="2"/>
      <text x="${padL}" y="${padT - 6}" font-size="11" fill="#6b7684">浓度 (g/m³, 对数轴)</text>
      <text x="${W - padR}" y="${H - 6}" text-anchor="end" font-size="11" fill="#6b7684">x (m)</text>
    </svg>`;
}

function currentRequest(isAxis) {
  const source = buildSource();
  if (isAxis) {
    return {
      q: fields.q(),
      source,
      wind_speed: fields.u(),
      stability: fields.stab(),
      axis: { start: fields.x0(), end: fields.x1(), count: fields.n(), y: fields.y(), z: fields.z() },
    };
  }
  return {
    q: fields.q(),
    source,
    wind_speed: fields.u(),
    stability: fields.stab(),
    receptor: { x: fields.x(), y: fields.y(), z: fields.z() },
  };
}

async function runConc() {
  setBusy(true);
  setError("");
  try {
    const p = await api("/api/conc", currentRequest(false));
    renderPoint(p);
  } catch (e) {
    setError(String(e.message || e));
  } finally {
    setBusy(false);
  }
}

async function runAxis() {
  setBusy(true);
  setError("");
  try {
    const a = await api("/api/axis", currentRequest(true));
    renderAxis(a);
  } catch (e) {
    setError(String(e.message || e));
  } finally {
    setBusy(false);
  }
}

async function loadExample() {
  setBusy(true);
  setError("");
  try {
    const res = await fetch("/example/ground-level.json");
    if (!res.ok) {
      throw new Error("示例文件加载失败：HTTP " + res.status);
    }
    const data = await res.json();
    $("f-q").value = data.q;
    $("f-u").value = data.wind_speed;
    $("f-stab").value = data.stability;
    $("f-x0").value = data.axis.start;
    $("f-x1").value = data.axis.end;
    $("f-n").value = data.axis.count;
    $("f-y").value = data.axis.y;
    $("f-z").value = data.axis.z;
    setMode("height");
    $("f-height").value = data.source.height;
    await runAxis();
  } catch (e) {
    setError(String(e.message || e));
  } finally {
    setBusy(false);
  }
}

function setMode(mode) {
  document.querySelectorAll("#sourceMode button").forEach((b) => {
    b.classList.toggle("active", b.dataset.mode === mode);
  });
  $("fields-height").style.display = mode === "height" ? "block" : "none";
  $("fields-stack").style.display = mode === "stack" ? "block" : "none";
}

document.querySelectorAll("#sourceMode button").forEach((b) => {
  b.addEventListener("click", () => setMode(b.dataset.mode));
});
$("btn-conc").addEventListener("click", runConc);
$("btn-axis").addEventListener("click", runAxis);
$("btn-example").addEventListener("click", loadExample);
