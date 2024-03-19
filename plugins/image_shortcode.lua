function HandleShortcode(args)
    if #args == 1 then
        local imageSrc = string.format("/media/%s", args[1])

        return string.format("![image](%s)", imageSrc)
    elseif #args == 2 then
        local imageSrc = string.format("/media/%s", args[1])

        return string.format("![%s](%s)", args[2], imageSrc)
    else
        return ""
    end
end
