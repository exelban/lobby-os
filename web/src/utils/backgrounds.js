export const backgrounds = [
  { id: "glow", label: "Color glow" },
  { id: "neutral", label: "Quiet neutral" },
  { id: "stone", label: "Warm stone" },
  { id: "mist", label: "Soft mist" },
  { id: "halo", label: "Single halo" },
  { id: "dots", label: "Precision dots" },
  { id: "sage", label: "Muted sage" },
  { id: "grid", label: "Fine grid" },
  { id: "dusk", label: "Dusty dusk" },
  { id: "contour", label: "Contour lines" },
  { id: "aurora", label: "Soft aurora" },
  { id: "legacy", label: "Legacy" },
]

export const normalizeBackground = value => backgrounds.some(({ id }) => id === value) ? value : "legacy"
