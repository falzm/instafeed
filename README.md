# Instafeed

`instafeed` is a very simple utility that allows you to generate a RSS feed of an [Instagram](https://www.instagram.com/) user.

## Requirements

* The Go language compiler (version >= 1.9)
* A valid Instagram account

## Building

At the top of the sources directory, just type `make`. If everything went well, you should end up with binary named `instafeed` in your current directory.

## Usage

`instafeed` expects the `IG_LOGIN` and `IG_PASSWORD` environment variables set to your Instagram login and password respectively, and the username of the Instagram user provided as argument. On successful execution, it prints the resulting RSS feed on the standard output.

Example:

```
$ export IG_LOGIN="your_instagram_login" IG_PASSWORD="********"
$ ./instafeed marutaro > marutaro.xml
$ cat marutaro.xml
<?xml version="1.0" encoding="UTF-8"?><rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <title>Instagram Feed: Shinjiro Ono</title>
    <link>https://www.instagram.com/marutaro</link>
    <description>My official shop is in 2-12-3 Nezu Bunkyo-ku Japan✨🐶✨ まるのお店は 根津駅徒歩1分にあるよ！👇ブログ更新したよ👇</description>
    <managingEditor>marutaro@instagram.com (Shinjiro Ono)</managingEditor>
    <pubDate>Tue, 03 Apr 2018 15:45:15 +0200</pubDate>
    <item>
      <title>You&#39;ve had a tough day. Thanks for your hard wo...</title>
      <link>https://www.instagram.com/p/BhGv7aMnkAm</link>
      <description>You&#39;ve had a tough day. Thanks for your hard wo...</description>
      <content:encoded><![CDATA[<p>You've had a tough day. Thanks for your hard work.✨🐶✨おつまる〜＼(^o^)／ 今日もよく頑張ったね！
#今日のお昼は何食べて来たの❓
#っていうか昨日のお昼何食べたか覚えてる❓</p><p><img src="https://scontent-cdg2-1.cdninstagram.com/vp/398dc348ba778a65a8313936c3b3b6b4/5B5A6F96/t51.2885-15/s750x750/sh0.08/e35/29417321_2040161612665341_6846638084159700992_n.jpg?ig_cache_key=MTc0OTI5NjI5NjA0NDE1MDgyMg%3D%3D.2"></p>]]></content:encoded>
      <pubDate>Tue, 03 Apr 2018 12:46:30 +0200</pubDate>
    </item>
    ...
  </channel>
</rss>
```
