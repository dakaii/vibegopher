#!/usr/bin/env node
/**
 * vue-tsc still resolves TypeScript's classic Node API (typescript/lib/tsc).
 * TypeScript 7 is native and does not export that path yet, so we bridge to
 * @typescript/typescript6 for Vue SFC typechecking while the project stays on TS 7.
 */
const Module = require('node:module')
const path = require('node:path')

const originalResolve = Module._resolveFilename
const ts6Root = path.dirname(require.resolve('@typescript/typescript6/package.json'))

Module._resolveFilename = function patchedResolve(request, parent, isMain, options) {
  if (request === 'typescript') {
    return path.join(ts6Root, 'lib', 'typescript.js')
  }
  if (request === 'typescript/lib/tsc' || request === 'typescript/lib/tsc.js') {
    return path.join(ts6Root, 'lib', 'tsc.js')
  }
  return originalResolve.call(this, request, parent, isMain, options)
}

process.argv.splice(2, 0, '--noEmit', '-p', 'tsconfig.app.json')
require('vue-tsc/index.js').run()
