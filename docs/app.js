/**
 * osx-traffic-stats — Interactive Landing Page Script
 * Geometric connected dots & lines canvas, screenshot showcase tabs, and clipboard.
 */

document.addEventListener('DOMContentLoaded', () => {
  initBackgroundCanvas();
  initShowcaseTabs();
  initInstallPill();
  initMiniCopyButtons();
});

/* ==========================================================================
   1. Subtle Geometric Connected Lines & Dots Canvas
   ========================================================================== */

function initBackgroundCanvas() {
  const canvas = document.getElementById('bg-canvas');
  if (!canvas) return;

  const ctx = canvas.getContext('2d');
  let width = (canvas.width = window.innerWidth);
  let height = (canvas.height = window.innerHeight);

  const points = [];
  const count = Math.min(Math.floor((width * height) / 18000), 75);
  const maxDistance = 140;

  for (let i = 0; i < count; i++) {
    points.push({
      x: Math.random() * width,
      y: Math.random() * height,
      vx: (Math.random() - 0.5) * 0.45,
      vy: (Math.random() - 0.5) * 0.45,
      radius: Math.random() * 1.5 + 1.0,
      color: Math.random() > 0.5 ? 'rgba(6, 182, 212, ' : 'rgba(16, 185, 129, '
    });
  }

  function resize() {
    width = canvas.width = window.innerWidth;
    height = canvas.height = window.innerHeight;
  }

  window.addEventListener('resize', resize, { passive: true });

  let mouse = { x: -1000, y: -1000 };
  window.addEventListener('mousemove', (e) => {
    mouse.x = e.clientX;
    mouse.y = e.clientY;
  }, { passive: true });

  function draw() {
    ctx.clearRect(0, 0, width, height);

    // Update & draw nodes
    for (let i = 0; i < points.length; i++) {
      const p = points[i];
      p.x += p.vx;
      p.y += p.vy;

      if (p.x < 0) p.x = width;
      if (p.x > width) p.x = 0;
      if (p.y < 0) p.y = height;
      if (p.y > height) p.y = 0;

      // Mouse gentle repulsion
      const dxm = p.x - mouse.x;
      const dym = p.y - mouse.y;
      const distM = Math.sqrt(dxm * dxm + dym * dym);
      if (distM < 100) {
        p.x += (dxm / distM) * 0.6;
        p.y += (dym / distM) * 0.6;
      }

      ctx.beginPath();
      ctx.arc(p.x, p.y, p.radius, 0, Math.PI * 2);
      ctx.fillStyle = p.color + '0.6)';
      ctx.fill();

      // Connect nearby nodes
      for (let j = i + 1; j < points.length; j++) {
        const p2 = points[j];
        const dx = p.x - p2.x;
        const dy = p.y - p2.y;
        const dist = Math.sqrt(dx * dx + dy * dy);

        if (dist < maxDistance) {
          const alpha = (1 - dist / maxDistance) * 0.18;
          ctx.beginPath();
          ctx.moveTo(p.x, p.y);
          ctx.lineTo(p2.x, p2.y);
          ctx.strokeStyle = `rgba(6, 182, 212, ${alpha})`;
          ctx.lineWidth = 0.85;
          ctx.stroke();
        }
      }
    }

    requestAnimationFrame(draw);
  }

  requestAnimationFrame(draw);
}

/* ==========================================================================
   2. Showcase View Switching Tabs
   ========================================================================== */

function initShowcaseTabs() {
  const tabs = document.querySelectorAll('.showcase-tab');
  const panels = document.querySelectorAll('.showcase-panel');

  tabs.forEach((tab) => {
    tab.addEventListener('click', () => {
      const targetView = tab.getAttribute('data-view');

      // Update tab active classes
      tabs.forEach((t) => {
        t.classList.remove('active');
        t.setAttribute('aria-selected', 'false');
      });
      tab.classList.add('active');
      tab.setAttribute('aria-selected', 'true');

      // Update visible panel
      panels.forEach((p) => p.classList.remove('active'));
      const activePanel = document.getElementById(`view-${targetView}`);
      if (activePanel) {
        activePanel.classList.add('active');
      }
    });
  });
}

/* ==========================================================================
   3. Pill-Shaped Install Input & Copy Action
   ========================================================================== */

function initInstallPill() {
  const input = document.getElementById('install-cmd-input');
  const btnCopy = document.getElementById('btn-copy-install');
  const toggleCask = document.getElementById('btn-toggle-cask');

  const cliCmd = 'brew install smford/tap/osx-traffic-stats';
  const caskCmd = 'brew install --cask smford/tap/osx-traffic-stats';

  if (btnCopy && input) {
    btnCopy.addEventListener('click', () => {
      copyToClipboard(input.value, btnCopy);
    });
  }

  if (toggleCask && input) {
    toggleCask.addEventListener('click', () => {
      if (input.value === caskCmd) {
        input.value = cliCmd;
        toggleCask.textContent = caskCmd;
      } else {
        input.value = caskCmd;
        toggleCask.textContent = cliCmd;
      }
    });
  }
}

function initMiniCopyButtons() {
  const miniButtons = document.querySelectorAll('.btn-mini-copy');
  miniButtons.forEach((btn) => {
    btn.addEventListener('click', () => {
      const text = btn.getAttribute('data-copy');
      if (text) {
        navigator.clipboard.writeText(text).then(() => {
          const original = btn.textContent;
          btn.textContent = '✓';
          setTimeout(() => {
            btn.textContent = original;
          }, 1500);
        });
      }
    });
  });
}

function copyToClipboard(text, btnElement) {
  navigator.clipboard.writeText(text).then(() => {
    const textSpan = btnElement.querySelector('.copy-text');
    if (textSpan) {
      const originalText = textSpan.textContent;
      textSpan.textContent = 'Copied!';
      btnElement.classList.add('copied');
      setTimeout(() => {
        textSpan.textContent = originalText;
        btnElement.classList.remove('copied');
      }, 2000);
    }
  });
}

