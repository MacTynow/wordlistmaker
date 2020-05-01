# Wordlistmaker

Google cloud function to translate a list of english and word in another language and returns a csv with english, target language and (optional) romanization.

## Usage

```json
# example.json
{
    "words": [
        "work",
        "a day",
        "朋友",
        "client",
        "你好"
    ]
}
```

returns

```csv
work, 工作, gōng zuò
a day, 一天, yī tiān
friend, 朋友, péng yǒu
client, 客户, kè hù
hello there, 你好, nǐ hǎo
```