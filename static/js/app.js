// Global UI enhancements
// - scroll reveal on elements with .sr-hidden
// - ripple effect on buttons/links
// - back-to-top button
// - smooth hover tilt on cards

(function () {
  const qs = (sel) => document.querySelector(sel);
  const qsa = (sel) => Array.from(document.querySelectorAll(sel));

  document.addEventListener('DOMContentLoaded', () => {
    // Back to top
    let backTop = qs('.back-to-top');
    if (!backTop) {
      backTop = document.createElement('button');
      backTop.className = 'back-to-top';
      backTop.textContent = '↑';
      document.body.appendChild(backTop);
    }
    backTop.addEventListener('click', () => {
      window.scrollTo({ top: 0, behavior: 'smooth' });
    });
    const onScrollTop = () => {
      const show = window.scrollY > 200;
      backTop.classList.toggle('show', show);
    };
    window.addEventListener('scroll', onScrollTop);
    onScrollTop();

    // Scroll reveal
    const revealEls = qsa('.fade-in, .slide-in-up, .sr-hidden');
    const onReveal = () => {
      const vh = window.innerHeight;
      revealEls.forEach((el, idx) => {
        const rect = el.getBoundingClientRect();
        if (rect.top < vh - 40) {
          el.style.transitionDelay = `${Math.min(idx * 50, 400)}ms`;
          el.classList.add('enter', 'sr-show');
          el.classList.remove('sr-hidden');
        }
      });
    };
    window.addEventListener('scroll', onReveal, { passive: true });
    onReveal();

    // Ripple effect
    const attachRipple = (el) => {
      el.addEventListener('click', (e) => {
        const rect = el.getBoundingClientRect();
        const ripple = document.createElement('span');
        ripple.className = 'ripple';
        const size = Math.max(rect.width, rect.height);
        ripple.style.width = ripple.style.height = `${size}px`;
        ripple.style.left = `${e.clientX - rect.left - size / 2}px`;
        ripple.style.top = `${e.clientY - rect.top - size / 2}px`;
        el.appendChild(ripple);
        ripple.addEventListener('animationend', () => ripple.remove());
      });
    };
    qsa('button, input[type="submit"], .nav-links a, .view-switcher button').forEach(attachRipple);

    // Hover tilt (subtle) for cards and list items
    const tiltTargets = qsa('.card, li');
    tiltTargets.forEach((el) => {
      el.addEventListener('mousemove', (e) => {
        const rect = el.getBoundingClientRect();
        const dx = (e.clientX - rect.left) / rect.width - 0.5;
        const dy = (e.clientY - rect.top) / rect.height - 0.5;
        el.style.transform = `perspective(700px) rotateX(${(-dy * 3)}deg) rotateY(${dx * 3}deg)`;
      });
      el.addEventListener('mouseleave', () => {
        el.style.transform = 'perspective(700px) rotateX(0) rotateY(0)';
      });
    });
  });
})();
