hostname=$(cat /etc/hostname)
source ~/.bashrc
source ~/.vimrc

echo
echo "#######################################"
echo " MONTANA STATE UNIVERSITY"
echo " NeuroCAVE"
echo " Connected to $hostname"
date
uptime
echo "#######################################"
echo

#Create Shortcuts
alias gu="git pull"
alias zu="echo \"Moving to msuclassfiles\"; cd ~/Desktop/msuclassfiles/; echo \"Pulling from git repository\"; git pull"
alias gst="git status"
alias ga="git add ."
alias gm="git merge"
alias gc="git commit -m"
alias gp="git push -u origin master"

alias ipp="curl ipecho.net/plain; echo;"
alias ipi="ifconfig | grep \"inet \" | grep -v 127.0.0.1"
