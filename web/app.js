async function jget(url) {
  const r = await fetch(url);
  return r.json();
}

async function jpost(url, body) {
  const r = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const text = await r.text();
  if (!r.ok) throw new Error(text || r.statusText);
  try { return JSON.parse(text); } catch { return text; }
}

async function refresh() {
  document.getElementById("stats").textContent = JSON.stringify(await jget("/api/stats"), null, 2);
  document.getElementById("fences").textContent = JSON.stringify(await jget("/api/fences"), null, 2);
  document.getElementById("alerts").textContent = JSON.stringify(await jget("/api/alerts"), null, 2);
}

document.getElementById("reg").onclick = async () => {
  const msg = document.getElementById("msg");
  try {
    const vertices = JSON.parse(document.getElementById("verts").value);
    await jpost("/api/fence", {
      id: document.getElementById("fid").value,
      name: document.getElementById("fname").value,
      kind: 0,
      vertices,
      alert_on_enter: true,
      alert_on_exit: true,
      alert_on_dwell: true,
      dwell_after: 0,
    });
    msg.textContent = "围栏已注册";
    await refresh();
  } catch (e) {
    msg.textContent = String(e);
  }
};

document.getElementById("ingest").onclick = async () => {
  const msg = document.getElementById("msg");
  try {
    await jpost("/api/ingest", {
      object_id: document.getElementById("oid").value,
      lat: parseFloat(document.getElementById("lat").value),
      lng: parseFloat(document.getElementById("lng").value),
      at: new Date().toISOString(),
    });
    msg.textContent = "已注入";
    await refresh();
  } catch (e) {
    msg.textContent = String(e);
  }
};

document.getElementById("refresh").onclick = () => refresh().catch(console.error);
refresh().catch(console.error);
setInterval(() => refresh().catch(() => {}), 4000);
