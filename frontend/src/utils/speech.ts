// Thin wrapper around the browser's Web Speech API for reading Chinese
// words/sentences aloud. No backend, no audio files — quality (and even
// availability) depends on whatever Chinese voice the user's OS/browser
// ships with. Phones (Android/iOS) almost always have one built in;
// desktop OSes (notably Windows) often don't unless the user installed the
// matching language pack, so this module tracks whether one is actually
// available and lets the UI reflect that instead of failing silently.
import { ref } from 'vue'

export function isSpeechSupported(): boolean {
  return typeof window !== 'undefined' && 'speechSynthesis' in window
}

const chineseVoiceAvailable = ref(false)

function refreshVoiceAvailability() {
  if (!isSpeechSupported()) {
    return
  }

  chineseVoiceAvailable.value = window.speechSynthesis
    .getVoices()
    .some(voice => voice.lang.toLowerCase().startsWith('zh'))
}

if (isSpeechSupported()) {
  refreshVoiceAvailability()
  window.speechSynthesis.addEventListener('voiceschanged', refreshVoiceAvailability)
}

// Reactive — true once a zh-* voice is found. May start false and flip to
// true shortly after load, since voice lists often populate asynchronously.
export function useChineseVoiceAvailable() {
  return chineseVoiceAvailable
}

export function speak(text: string, lang = 'zh-CN'): void {
  if (!isSpeechSupported() || !text) {
    return
  }

  // Cancel whatever's currently playing so rapid clicks don't queue up.
  window.speechSynthesis.cancel()

  const utterance = new SpeechSynthesisUtterance(text)

  utterance.lang = lang
  utterance.rate = 0.85

  window.speechSynthesis.speak(utterance)
}
