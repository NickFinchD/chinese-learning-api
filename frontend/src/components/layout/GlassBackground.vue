<template>
  <div class="fixed inset-0 -z-10 overflow-hidden bg-gradient-to-b from-gray-100 via-gray-200 to-gray-300 transition-colors duration-300 dark:from-slate-950 dark:via-slate-950 dark:to-slate-900">
    <div class="animate-blob absolute -left-40 -top-40 h-[28rem] w-[28rem] rounded-full bg-[var(--color-primary)] opacity-30 blur-3xl dark:opacity-25" />
    <div
      class="animate-blob absolute -right-32 top-0 h-[26rem] w-[26rem] rounded-full bg-[var(--color-mint)] opacity-25 blur-3xl dark:opacity-20"
      style="animation-delay: -6s; animation-duration: 24s"
    />

    <!-- Layered mountain skyline, recoloured from the active palette via the
         --color-* custom properties (see stores/palette.ts) so it matches
         whatever gamma the user picks. Blurred so the glass cards' own
         backdrop-blur has soft, painterly colour to pick up rather than a
         crisp graphic. -->
    <svg
      class="absolute inset-x-0 bottom-0 h-2/3 w-full min-w-[1440px] blur-sm"
      viewBox="0 0 1440 480"
      preserveAspectRatio="xMidYMax slice"
      aria-hidden="true"
    >
      <path
        d="M0,480 L0,300 L120,255 L240,300 L360,220 L480,275 L600,235 L720,300 L840,250 L960,305 L1080,260 L1200,290 L1320,240 L1440,280 L1440,480 Z"
        fill="var(--color-mint)"
        opacity="0.35"
      />
      <path
        d="M0,480 L0,340 L100,260 L220,320 L340,210 L460,300 L580,190 L700,280 L820,170 L940,290 L1060,220 L1180,300 L1300,230 L1440,270 L1440,480 Z"
        fill="var(--color-accent)"
        opacity="0.35"
      />
      <path
        d="M0,480 L0,380 L150,250 L300,350 L450,190 L600,330 L750,160 L900,320 L1050,200 L1200,340 L1350,220 L1440,290 L1440,480 Z"
        fill="var(--color-primary)"
        opacity="0.45"
      />
    </svg>
  </div>
</template>

<script setup lang="ts">
</script>
