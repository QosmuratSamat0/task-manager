// Intersection-based reveal animations and small polish
(function(){
  const supportsIO = 'IntersectionObserver' in window;
  const prefersReduced = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;

  function applyReveal() {
    const nodes = Array.from(document.querySelectorAll('.reveal'));
    if (prefersReduced || !supportsIO) {
      nodes.forEach(n => n.classList.add('visible'));
      return;
    }
    const io = new IntersectionObserver((entries) => {
      entries.forEach(e => {
        if (e.isIntersecting) {
          const delay = parseInt(e.target.getAttribute('data-reveal-delay') || '0', 10);
          setTimeout(() => e.target.classList.add('visible'), Math.max(0, delay));
          io.unobserve(e.target);
        }
      });
    }, { rootMargin: '0px 0px -10% 0px', threshold: 0.15 });
    nodes.forEach(n => io.observe(n));
  }

  // Run after DOM painted
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', applyReveal);
  } else {
    applyReveal();
  }
})();

