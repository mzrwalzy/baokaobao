#!/bin/bash
# 同步推送脚本 - 同时推送到 GitHub 和微信 GitLab

cd /home/project/baokaobao

echo "正在推送到 GitHub..."
git push origin HEAD:master

echo "正在推送到微信 GitLab..."
git push wxgit HEAD:master

echo "推送完成！"
