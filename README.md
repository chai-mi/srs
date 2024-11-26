# Domain generate
一个域名分类方案，使用 `tag` 和 `weight` 标记域名

## 如何影响输出
在 `extra.json` 中为域名添加 `tag`，`tag` 会为 `config.json` 中对应 `tag` 以 `rules.tag weight` 的值增添权重，权重大于 0 的域名即属于该类