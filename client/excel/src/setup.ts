// Setup file for vitest to mock Office.js before any tests run.
import { vi } from 'vitest'

// Mock settings.
const mockSettings: Record<string, string> = {
  KIBTAB_HOST: 'http://localhost:8080',
  KIBTAB_TABLE: 'users',
}

// Mock the Office object.
const mockOffice = {
  context: {
    document: {
      setSelectedDataAsync: vi.fn(),
      onSelectionChanged: null as ((handler: () => void) => void) | null,
      reload: vi.fn(),
    },
    settings: {
      get: vi.fn((key: string) => mockSettings[key]),
    },
  },
}

// @ts-ignore - Mock the global Office object.
;(global as any).Office = mockOffice