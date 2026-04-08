#!/bin/bash
cd /home/project/baokaobao

echo "正在推送到 GitHub..."
git push origin HEAD:master

echo "正在推送到微信 GitLab..."
git push wxgit HEAD:master

echo "推送完成！"
