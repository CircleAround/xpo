const jaPlaceholderMarkdown = `# Markdown で書けます
\`Markdown\` のサンプル

## 中タイトル
### 小タイトル

- 順序関係ないリスト 1
- 順序関係ないリスト 2
- 順序関係ないリスト 3


1. 順序関係あるリスト 1
1. 順序関係あるリスト 2
1. 順序関係あるリスト 3

\`\`\`
ソースコードとか
\`\`\`
`

const jaHelpMarkdown = `
テキストエリアにマークダウンで記述できます。
変更の反映はショートカットキーでもできます。
`

export default {
  en: {
    message: {
      hello: 'hello world'
    }
  },
  ja: {
    message: {
      hello: 'こんにちは、世界'
    },
    ui: {
      placeholder: {
        markdown: jaPlaceholderMarkdown
      },
      help: {
        markdown: jaHelpMarkdown,
        shortcutkey: {
          post: {
            mac: 'Command + Enter',
            win: 'Control + Enter'
          }
        }
      }
    }
  }
}
