{
  // MelodyMe's conductor page. The conductor watches the pointer, and goes cross-eyed and dizzy when
  // it is between the eyes.
  const pupils = [...document.querySelectorAll("[data-watch] .cd-pupil")];
  if (pupils.length === 2 && !matchMedia("(prefers-reduced-motion: reduce)").matches) {
    const conductor = pupils[0].closest("[data-watch]");
    let dizzy = false;
    let eyes; // where the eyes are: centres and width, which a blink doesn't change
    const look = (p, x, y) => (p.style.transform = `translate(${x.toFixed(1)}px,${y.toFixed(1)}px)`);
    addEventListener("pointermove", (e) => {
      const x = e.clientX, y = e.clientY;
      // Measured afresh only while not dizzy: the sway would carry the zone away from the pointer.
      if (!dizzy || !eyes) {
        eyes = pupils.map((p) => {
          const b = p.previousElementSibling.getBoundingClientRect();
          return { x: b.left + b.width / 2, y: b.top + b.height / 2, w: b.width };
        });
      }
      const [l, r] = eyes;
      // Dizzy starts only in the gap between the eyes, level with them, and ends only once the
      // pointer is past the eyes' outer edges, so passing slowly over the edge doesn't flick it.
      if (dizzy) {
        dizzy = x > l.x - l.w / 2 && x < r.x + r.w / 2 && y > l.y - l.w && y < l.y + l.w;
      } else {
        dizzy = x > l.x + l.w / 2 && x < r.x - r.w / 2 && y > l.y - l.w / 2 && y < l.y + l.w / 2;
      }
      conductor.classList.toggle("is-dizzy", dizzy);
      if (dizzy) {
        look(pupils[0], 7.5, 2.5);
        look(pupils[1], -7.5, 2.5);
        return;
      }
      pupils.forEach((p, i) => {
        const eye = i ? r : l;
        const dx = x - eye.x;
        const dy = y - eye.y;
        const d = Math.hypot(dx, dy) || 1;
        const m = Math.min(6, d / 25);
        look(p, (dx / d) * m, (dy / d) * m);
      });
    });
  }
}
