## Encode / decode base 64

<details>
<summary>Answer</summary>

```

echo -n "string" | base64 -w 0
echo -n "string" | base64 -d

```

-n is important to avoid new line character
- w 0 is important for long multi-line tokens

</details>

