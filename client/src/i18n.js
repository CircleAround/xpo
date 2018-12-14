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
    error: {
      messages: {
        validation: {
          required: `{property}は必須です`,
          min: `{property}は短すぎます。`,
          max: `{property}は長すぎます。`,
          identity_name_format: `{property}は先頭英数小文字かつ半角英数小文字とアンダースコアです`,
          usernickname_format: `ニックネームには <>/:"'と空白を含めてはいけません`,
          nothing: `{property}が何らかのエラーです`
        },
        valueNotUnique: `{property}は既に存在します`,
        invalidParameter: `{property}は利用できません`,
        duplicatedObject: `既に存在します`,
        duplicatedUser: 'すでにユーザー登録済みです',
        duplicatedUserName:
          '既に取得されてしまったユーザー名です。別の名前にしましょう。',
        unexpected: '予期しないエラーです'
      }
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
