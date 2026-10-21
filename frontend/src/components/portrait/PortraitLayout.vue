<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { gsap } from 'gsap'
import portrait from '../../assets/portrait-gold.jpeg'
import './portrait-board.css'

// 从 ui-portrait-music 提取布局与入场动效；业务内容通过 slot 传入。
const props = defineProps({
  duration: { type: Number, default: 640 },
  reduced: { type: Boolean, default: false },
  paused: { type: Boolean, default: false },
})
const root = ref(null)
let media, panelTween, backgroundTween

onMounted(() => {
  media = gsap.matchMedia()
  // matchMedia 自带 context；选择器只在当前组件内生效。
  media.add('(prefers-reduced-motion: no-preference)', () => {
    if (props.reduced || props.paused) return
    const duration = Number.isFinite(props.duration) ? Math.min(1600, Math.max(120, props.duration)) / 1000 : 0.64
    panelTween = gsap.from('[data-motion-layer]', {
      autoAlpha: 0, x: -20, duration, stagger: 0.07,
      ease: 'power3.out', clearProps: 'transform,opacity,visibility',
    })
    backgroundTween = gsap.from('.pm-backdrops img', {
      scale: 1.035, duration: 0.85, ease: 'power3.out', clearProps: 'transform',
    })
  }, root.value)
})

watch(() => props.reduced || props.paused, stopped => {
  if (stopped) {
    panelTween?.progress(1)
    backgroundTween?.progress(1)
  }
})
onBeforeUnmount(() => media?.revert())
</script>

<template>
  <div ref="root" class="pm-room pm-board">
    <div class="pm-backdrops" aria-hidden="true">
      <img :src="portrait" alt="" fetchpriority="high" />
    </div>
    <header class="pm-masthead">
      <a href="#publish" class="pm-brand"><h1>校园留言板</h1><span aria-hidden="true">↗</span></a>
    </header>
    <main class="pm-column">
      <nav class="pm-tabs" aria-label="留言板导航">
        <a href="#publish"><span>发布留言</span><span aria-hidden="true">↗</span></a>
        <a href="#messages"><span>查看留言</span><span aria-hidden="true">↓</span></a>
      </nav>
      <slot />
    </main>
  </div>
</template>
